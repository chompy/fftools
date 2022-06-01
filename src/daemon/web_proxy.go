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
	"io"
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
	// build buffers
	msgTypeBuf := make([]byte, 1)
	reqInfoBuf := make([]byte, 4)
	// itterate and read connection
	for {

		if _, err := io.ReadFull(conn, msgTypeBuf); err != nil {
			logWarn("[PROXY] " + err.Error())
			return err
		}
		switch msgTypeBuf[0] {
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
				// read message info in to buffer
				if _, err := io.ReadFull(conn, reqInfoBuf); err != nil {
					logWarn("[PROXY] " + err.Error())
					return err
				}
				rid := binary.LittleEndian.Uint16(reqInfoBuf[0:])
				reqLen := binary.LittleEndian.Uint16(reqInfoBuf[2:])
				// read remaining message
				reqPathBuf := make([]byte, reqLen)
				if _, err := io.ReadFull(conn, reqPathBuf); err != nil {
					logWarn("[PROXY] " + err.Error())
					return err
				}
				reqPath := string(reqPathBuf)
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
				out := make([]byte, 7)
				out[0] = proxyMsgWebResp
				binary.LittleEndian.PutUint16(out[1:], rid)
				binary.LittleEndian.PutUint32(out[3:], uint32(buf.Len()))
				if buf.Len()+7 > proxyResponseMaxSize {
					return ErrProxyResponseTooLarge
				}
				out = append(out, buf.Bytes()...)
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
