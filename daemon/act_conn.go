package main

import (
	"net"
)

const dataTypeActScripts = 201
const dataTypeActScriptEnable = 202
const dataTypeActScriptDisable = 203
const dataTypeActScriptReload = 204
const dataTypeActPlayer = 205
const dataTypeActSay = 206
const dataTypeActEnd = 207

var actConn *net.UDPConn = nil
var remoteAddr *net.UDPAddr = nil
var hasRequestedPlayer = false

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
		// once connection is established request that act send last player change line
		if !hasRequestedPlayer {
			hasRequestedPlayer = true
			actRequestPlayer()
		}
		// decode
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
		}

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
	for _, script := range scripts {
		enabledString := ""
		if script.Enabled {
			enabledString = "1"
		}
		lastErrMsg := ""
		if script.LastError != nil {
			lastErrMsg = script.LastError.Error()
		}
		data := script.ScriptName + "|" + enabledString + "|" + script.Name + "|" + script.Desc + "|" + lastErrMsg
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

func actRequestPlayer() error {
	logInfo("[ACT] Request primary player.")
	return actRawSend([]byte{byte(dataTypeActPlayer)})
}
