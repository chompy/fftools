package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
)

const webProxyAddress = "localhost:31595"

func webProxyConnect() error {
	conn, err := net.Dial("tcp", webProxyAddress)
	if err != nil {
		logWarn(err.Error())
		return err
	}
	defer conn.Close()
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
			rid := binary.LittleEndian.Uint16(buf)
			reqLen := binary.LittleEndian.Uint16(buf[2:])
			reqPath := string(buf[4 : reqLen+4])
			logInfo("[PROXY] GET %s", reqPath)
			// handle request
			w := newWebProxyResponseWriter()
			r, err := http.NewRequest(http.MethodGet, "http://localhost"+reqPath, nil)
			if err != nil {
				logWarn("[PROXY] " + err.Error())
				webServeB64(webError, http.StatusInternalServerError, w)
				continue
			}
			webHandle(w, r)
			// prepare response
			buf := bytes.Buffer{}
			w.r.Write(&buf)
			out := make([]byte, 6)
			binary.LittleEndian.PutUint16(out, rid)
			binary.LittleEndian.PutUint32(out[2:], uint32(buf.Len()))
			out = append(out, buf.Bytes()...)
			// send
			if _, err := conn.Write(out); err != nil {
				logWarn("[PROXY] " + err.Error())
				if errors.Is(err, io.EOF) {
					return err
				}
			}
		}
	}
}
