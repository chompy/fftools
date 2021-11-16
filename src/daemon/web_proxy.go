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
	"bytes"
	"encoding/binary"
	"net"
	"net/http"
)

const proxyMsgLogin = 1
const proxyMsgWebReq = 2
const proxyMsgWebResp = 3
const proxyMsgInvalidCreds = 4
const proxyResponseMaxSize = 32768

func webProxyConnect() error {
	config := configAppLoad()
	logInfo("[PROXY] Connecting to proxy server at %s.", config.ProxyAddress)
	conn, err := net.Dial("tcp", config.ProxyAddress)
	if err != nil {
		logWarn("[PROXY] " + err.Error())
		return err
	}
	defer conn.Close()
	uid, secret, err := configGetProxyCreds()
	logInfo("[PROXY] UID is %s.", uid)
	if err != nil {
		logWarn("[PROXY] " + err.Error())
		return err
	}
	if err := webProxySendCreds(uid, secret, conn); err != nil {
		logWarn("[PROXY] " + err.Error())
		return err
	}
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			logWarn("[PROXY] " + err.Error())
			return err
		}
		if n > 0 {
			switch buf[0] {
			case proxyMsgInvalidCreds:
				{
					// regenerate and resend creds if invalid
					// (assume collision)
					logInfo("Recieved invalid proxy creds, regenerating.")
					uid, secret := webProxyGenerateCreds()
					if err := configSetProxyCred(uid, secret); err != nil {
						logWarn("[PROXY] " + err.Error())
						return err
					}
					if err := webProxySendCreds(uid, secret, conn); err != nil {
						logWarn("[PROXY] " + err.Error())
						return err
					}
					break
				}
			case proxyMsgWebReq:
				{
					rid := binary.LittleEndian.Uint16(buf[1:])
					reqLen := binary.LittleEndian.Uint16(buf[3:])
					reqPath := string(buf[5 : reqLen+5])
					logInfo("[PROXY] #%d GET %s", rid, reqPath)
					// handle request
					w := newWebProxyResponseWriter()
					r, err := http.NewRequest(http.MethodGet, "http://localhost"+reqPath, nil)
					if err != nil {
						logWarn("[PROXY] " + err.Error())
						webServeB64(assetError500General, http.StatusInternalServerError, w)
						continue
					}
					webHandle(w, r)
					// prepare response
					buf := bytes.Buffer{}
					w.r.Write(&buf)
					out := make([]byte, proxyResponseMaxSize)
					out[0] = proxyMsgWebResp
					binary.LittleEndian.PutUint16(out[1:], rid)
					binary.LittleEndian.PutUint32(out[3:], uint32(buf.Len()))
					if buf.Len()+7 > proxyResponseMaxSize {
						return ErrProxyResponseTooLarge
					}
					bufBytes := buf.Bytes()
					for i := range bufBytes {
						out[i+7] = bufBytes[i]
					}
					// send
					if _, err := conn.Write(out); err != nil {
						logWarn("[PROXY] " + err.Error())
						return err
					}
					break
				}
			}

		}
	}
}
