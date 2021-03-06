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
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const proxyUserTTL = 86400
const maxCredWaitTime = 5000

type ProxyUser struct {
	Uid             string `json:"uid"`
	Secret          string `json:"secret"`
	connection      net.Conn
	lastRequestTime time.Time
	requests        []userRequest
	lastRequestId   uint16
	sync            sync.Mutex
}

var proxyUsers = make([]*ProxyUser, 0)

func addProxyUser(uid string, secret string, conn net.Conn) *ProxyUser {
	// uid and secret can't be empty
	if uid == "" || secret == "" {
		return nil
	}
	// fetch existing
	u := getProxyUser(uid)
	// create new
	if u == nil {
		u = &ProxyUser{
			Uid:    uid,
			Secret: secret,
		}
		proxyUsers = append(proxyUsers, u)
		if err := persistUsers(); err != nil {
			log.Printf("[WARN] persist users :: %s", err.Error())
		}
	}
	if u.Uid == "" || u.Secret == "" || u.Secret != secret {
		return nil
	}
	u.connection = conn
	u.requests = make([]userRequest, 0)
	u.lastRequestTime = time.Now()
	go func(u *ProxyUser) {
		msgTypeBuf := make([]byte, 1)
		for {
			// check that connection is set
			if u.connection == nil {
				log.Printf("[WARN] Nil connection found.")
				for i := range proxyUsers {
					if proxyUsers[i].Uid == u.Uid {
						proxyUsers = append(proxyUsers[i:], proxyUsers[:i+1]...)
						break
					}
				}
				return
			}
			// check response
			n, err := io.ReadFull(u.connection, msgTypeBuf)
			if err != nil {
				log.Printf("[WARN] proxy read response :: %s", err.Error())
				u.lastRequestTime = time.Time{}
			}
			// handle response
			if n > 0 {
				switch msgTypeBuf[0] {
				case proxyMsgWebResp:
					{
						if err := u.handleResponse(); err != nil {
							log.Printf("[WARN] proxy handle response :: %s", err.Error())
						}
						break
					}
				}
			}
			// clean up
			if time.Since(u.lastRequestTime) > time.Second*proxyUserTTL {
				log.Printf("[INFO] Clean up %s (%s).", u.connection.RemoteAddr().String(), u.Uid)
				u.connection.Close()
				for i := range proxyUsers {
					if proxyUsers[i].Uid == u.Uid {
						proxyUsers = append(proxyUsers[i:], proxyUsers[:i+1]...)
						break
					}
				}
				return
			}
		}
	}(u)
	return u
}

func getProxyUser(uid string) *ProxyUser {
	for _, u := range proxyUsers {
		if u.Uid == uid {
			return u
		}
	}
	return nil
}

func waitForCreds(conn net.Conn) (string, string) {
	start := time.Now()
	buf := make([]byte, 1+proxyUidLen+64)
	for {
		if time.Since(start) > time.Millisecond*maxCredWaitTime {
			conn.Close()
			log.Printf("[WARN] wait for uid, reached match wait time")
			return "", ""
		}
		_, err := conn.Read(buf)
		if err != nil {
			log.Printf("[WARN] wait for uid, read error :: %s", err.Error())
			conn.Close()
			return "", ""
		}
		if buf[0] != proxyMsgLogin {
			continue
		}
		uid := string(buf[1 : 1+proxyUidLen])
		secret := string(buf[10 : 1+proxyUidLen+64])
		return uid, secret
	}
}
