package main

import (
	"net"
	"strings"
)

const dataTypeActScriptsAvailable = 201
const dataTypeActScriptsEnabled = 202
const dataTypeActSay = 203
const dataTypeActEnd = 204
const actListenPort = 31593

var actConn *net.UDPConn = nil
var remoteAddr *net.UDPAddr = nil

func actListenUDP() error {
	addr := net.UDPAddr{
		Port: actListenPort,
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
		messageType := buf[0]
		switch messageType {
		case DataTypeLogLine:
			{
				// decode raw log + dispatch event
				logLine := LogLine{}
				if err := logLine.FromBytes(buf[:]); err != nil {
					break
				}
				eventListenerDispatch("act:log_line", logLine)
				// parsed log
				parsedLogEvent, err := ParseLogEvent(logLine)
				if err != nil {
					break
				}
				eventListenerDispatch("act:parsed_log_event", parsedLogEvent)
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

func actSendScriptsAvailable() error {
	scripts := luaGetAvailableScripts()
	data := strings.Join(scripts, ",")
	return actRawSend(
		append([]byte{byte(dataTypeActScriptsAvailable)}, []byte(data)...),
	)
}

func actSendScriptsEnabled() error {
	scripts := luaGetEnabledScripts()
	data := strings.Join(scripts, ",")
	return actRawSend(
		append([]byte{byte(dataTypeActScriptsEnabled)}, []byte(data)...),
	)
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
