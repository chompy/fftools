package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
)

const proxyMsgLogin = 1
const proxyMsgWebReq = 2
const proxyMsgWebResp = 3
const proxyMsgInvalidCreds = 4

func webProxyConnect() error {
	config := configAppLoad()
	logInfo("Connecting to proxy server at %s.", config.ProxyAddress)
	conn, err := net.Dial("tcp", config.ProxyAddress)
	if err != nil {
		logWarn(err.Error())
		return err
	}
	defer conn.Close()
	uid, secret, err := configGetProxyCreds()
	logInfo("Proxy UID is %s.", uid)
	if err != nil {
		logWarn(err.Error())
		return err
	}
	if err := webProxySendCreds(uid, secret, conn); err != nil {
		logWarn(err.Error())
		return err
	}
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			logWarn(err.Error())
			if errors.Is(err, io.EOF) {
				return err
			}
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
						logWarn(err.Error())
						return err
					}
					if err := webProxySendCreds(uid, secret, conn); err != nil {
						logWarn(err.Error())
						return err
					}
					break
				}
			case proxyMsgWebReq:
				{
					rid := binary.LittleEndian.Uint16(buf[1:])
					reqLen := binary.LittleEndian.Uint16(buf[3:])
					reqPath := string(buf[5 : reqLen+5])
					logInfo("[PROXY] GET %s", reqPath)
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
					out = append(out, buf.Bytes()...)
					// send
					if _, err := conn.Write(out); err != nil {
						logWarn("[PROXY] " + err.Error())
						if errors.Is(err, io.EOF) {
							return err
						}
					}
					break
				}
			}

		}
	}
}