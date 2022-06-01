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
	"log"
	"net"
)

const proxyPort = 31595
const proxyMsgLogin = 1
const proxyMsgWebReq = 2
const proxyMsgWebResp = 3
const proxyMsgInvalidCreds = 4
const proxyUidLen = 8

func proxyListen() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", proxyPort))
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("[WARN] proxyListen :: %s", err.Error())
			continue
		}
		log.Printf("[INFO] Connection from %s.", conn.RemoteAddr().String())
		go func(conn net.Conn) {
			uid, secret := waitForCreds(conn)
			if addProxyUser(uid, secret, conn) == nil {
				conn.Write([]byte{byte(proxyMsgInvalidCreds)})
				return
			}
		}(conn)
	}
}
