/*
This file is part of FFTools.

FFTools is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

FFTools is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with FFTools.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"net"
)

const dataTypeActScripts = 201
const dataTypeActScriptEnable = 202
const dataTypeActScriptDisable = 203
const dataTypeActScriptReload = 204
const dataTypeActScriptVersion = 205
const dataTypeActScriptUpdate = 206
const dataTypeActSay = 210
const dataTypeActEnd = 211
const dataTypeActUpdated = 212

var actConn *net.UDPConn = nil
var remoteAddr *net.UDPAddr = nil

func actListenUDP() error {
	config := configAppLoad()
	addr := net.UDPAddr{
		Port: int(config.PortData),
		IP:   net.ParseIP("127.0.0.1"),
	}
	var err error
	actConn, err = net.ListenUDP("udp", &addr)
	if err != nil {
		return err
	}
	defer actConn.Close()
	var buf [1024]byte
	for {
		rlen, remote, err := actConn.ReadFromUDP(buf[:])
		if err != nil || rlen == 0 {
			continue
		}
		remoteAddr = remote
		// decode
		go func(buf [1024]byte) {
			messageType := buf[0]
			switch messageType {
			case DataTypeLogLine:
				{
					// decode raw log + dispatch event
					logLine := LogLine{}
					if err := logLine.FromBytes(buf[:]); err != nil {
						break
					}
					// parsed log
					parsedLogEvent, err := ParseLogEvent(logLine)
					if err != nil {
						logWarn(err.Error())
						break
					}
					eventListenerDispatch("act:log_line", parsedLogEvent)
					break
				}
			case DataTypeCombatant:
				{
					combatant := Combatant{}
					if err := combatant.FromBytes(buf[:]); err != nil {
						break
					}
					eventListenerDispatch("act:combatant", combatant)
					break
				}
			case DataTypeEncounter:
				{
					encounter := Encounter{}
					if err := encounter.FromBytes(buf[:]); err != nil {
						break
					}
					eventListenerDispatch("act:encounter", encounter)
					break
				}
			case dataTypeActScripts:
				{
					if err := actSendScripts(); err != nil {
						logWarn(err.Error())
					}
					break
				}
			case dataTypeActScriptEnable, dataTypeActScriptDisable:
				{
					pos := 1
					scriptName := readString(buf[:], &pos)
					if messageType == dataTypeActScriptEnable {
						logInfo("[ACT] Enable '%s' script.", scriptName)
					} else {
						logInfo("[ACT] Disable '%s' script.", scriptName)
					}
					if err := configSetScriptEnabled(scriptName, messageType == dataTypeActScriptEnable); err != nil {
						logWarn(err.Error())
						break
					}
					luaEnableScripts()
					actSendScripts()
					break
				}
			case dataTypeActScriptReload:
				{
					pos := 1
					scriptName := readString(buf[:], &pos)
					for _, script := range loadedScripts {
						if script.ScriptName == scriptName {
							logInfo("[ACT] Reload '%s' script.", scriptName)
							if err := script.reload(); err != nil {
								logLuaWarn(script.L, err.Error())
							}
							break
						}
					}
					break
				}
			case dataTypeActScriptVersion:
				{
					pos := 1
					scriptName := readString(buf[:], &pos)
					if err := actScriptVersion(scriptName); err != nil {
						logWarn("[ACT] %s", err.Error())
					}
					break
				}
			case dataTypeActScriptUpdate:
				{
					pos := 1
					scriptName := readString(buf[:], &pos)
					if err := actScriptUpdate(scriptName); err != nil {
						logWarn("[ACT] %s", err.Error())
					}
					break
				}
			}
		}(buf)
	}
}

func actRawSend(data []byte) error {
	if actConn == nil || remoteAddr == nil {
		return ErrActNotConnected
	}
	if _, err := actConn.WriteTo(data, remoteAddr); err != nil {
		return err
	}
	return nil
}

func actSendScripts() error {
	logInfo("[ACT] Send script list.")
	scripts := luaLoadScripts()
	config := configAppLoad()
	uid, _, _ := configGetProxyCreds()
	webUrl := fmt.Sprintf("http://localhost:%d/", config.PortWeb)
	if config.EnableProxy {
		webUrl = fmt.Sprintf("%s%s/", config.ProxyURL, uid)
	}
	for _, script := range scripts {
		enabledString := ""
		if script.State == LuaScriptActive {
			enabledString = "1"
		}
		lastErrMsg := ""
		if script.LastError != nil {
			lastErrMsg = script.LastError.Error()
		}
		ver, _ := script.getCurrentVersion()
		verStr := "?"
		if ver != nil {
			verStr = ver.String()
		}
		data := script.ScriptName + "|" + enabledString + "|" + script.Name + "|" + script.Desc + "|" + verStr + "|" + webUrl + script.ScriptName + "|" + lastErrMsg
		err := actRawSend(append(
			[]byte{byte(dataTypeActScripts)},
			[]byte(data)...,
		))
		if err != nil {
			return err
		}
	}
	return nil
}

func actSay(text string) error {
	config := configAppLoad()
	if !config.EnableTTS {
		return nil
	}
	return actRawSend(
		append([]byte{byte(dataTypeActSay)}, []byte(text)...),
	)
}

func actEnd() error {
	return actRawSend(
		[]byte{byte(dataTypeActEnd)},
	)
}

func actError(err error, scriptName string) error {
	for _, script := range loadedScripts {
		if script.ScriptName == scriptName {
			script.LastError = err
			break
		}
	}
	return actSendScripts()
}

func actRequestUpdate() error {
	logInfo("[ACT] Request version update.")
	return actRawSend([]byte{byte(dataTypeActUpdated)})
}

func actScriptVersion(scriptName string) error {
	for _, script := range loadedScripts {
		if script.ScriptName == scriptName {
			logInfo("[ACT] Check '%s' script version.", scriptName)
			curVer, _ := script.getCurrentVersion()
			curVerStr := "?"
			if curVer != nil {
				curVerStr = curVer.String()
			}
			latestVer, _ := script.getLatestVersion()
			latestVerStr := "?"
			if latestVer != nil {
				latestVerStr = latestVer.String()
			}
			data := script.ScriptName + "|" + curVerStr + "|" + latestVerStr
			return actRawSend(append(
				[]byte{byte(dataTypeActScriptVersion)},
				[]byte(data)...,
			))
		}
	}
	return ErrLuaScriptNotFound
}

func actScriptUpdate(scriptName string) error {
	for _, script := range loadedScripts {
		if script.ScriptName == scriptName {
			logInfo("[ACT] Update '%s' script to latest version.", scriptName)
			if err := script.update(); err != nil {
				data := script.ScriptName + "|fail|" + err.Error()
				actRawSend(append(
					[]byte{byte(dataTypeActScriptUpdate)},
					[]byte(data)...,
				))
				return err
			}
			if err := script.reload(); err != nil {
				data := script.ScriptName + "|fail|" + err.Error()
				actRawSend(append(
					[]byte{byte(dataTypeActScriptUpdate)},
					[]byte(data)...,
				))
				return err
			}
			data := script.ScriptName + "|success"
			return actRawSend(append(
				[]byte{byte(dataTypeActScriptUpdate)},
				[]byte(data)...,
			))
		}
	}
	return ErrLuaScriptNotFound
}
