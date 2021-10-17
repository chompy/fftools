package main

import (
	"log"
	"net"
)

const listenPort = 31593

// ListenUDP listens for incomming data from ACT.
func ListenUDP() error {
	addr := net.UDPAddr{
		Port: listenPort,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	var buf [1024]byte
	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])
		if err != nil || rlen == 0 {
			continue
		}
		log.Println(rlen, remote, int(buf[0]))
		// decode
		messageType := buf[0]
		switch messageType {
		case DataTypeLogLine:
			{
				logLine := LogLine{}
				if err := logLine.FromBytes(buf[:]); err != nil {
					break
				}
				eventListenerDispatch("act:log_line", logLine)
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
		}

	}

}
