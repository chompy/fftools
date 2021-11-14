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
	"encoding/binary"
	"log"
	"net/http"
	"strings"
	"time"
)

const requestTimeout = 10000
const requestMaxSize = 2048
const responseMaxSize = 32768

type userRequest struct {
	id   uint16
	url  string
	resp []byte
}

func (u *ProxyUser) handleRequest(r *http.Request) (uint16, error) {
	defer u.sync.Unlock()
	u.sync.Lock()
	u.lastRequestTime = time.Now()
	u.lastRequestId++
	reqPath := "/" + strings.Join(strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")[1:], "/")
	ureq := userRequest{
		id:   u.lastRequestId,
		url:  reqPath + "?" + r.URL.RawQuery,
		resp: nil,
	}
	log.Printf("[INFO] [%s] %s %s", u.Uid, r.Method, ureq.url)
	u.requests = append(u.requests, ureq)
	out := make([]byte, 5)
	out[0] = proxyMsgWebReq
	binary.LittleEndian.PutUint16(out[1:], ureq.id)
	binary.LittleEndian.PutUint16(out[3:], uint16(len(ureq.url)))
	out = append(out, []byte(ureq.url)...)
	if len(out) > requestMaxSize {
		return 0, ErrRequestTooLarge
	}
	if _, err := u.connection.Write(out); err != nil {
		return 0, err
	}
	return u.lastRequestId, nil
}

func (u *ProxyUser) handleResponse(data []byte) error {
	if data[0] != proxyMsgWebResp {
		return ErrUnexpectedMessageType
	}
	id := binary.LittleEndian.Uint16(data[1:])
	respLen := binary.LittleEndian.Uint32(data[3:])
	if respLen+6 > responseMaxSize {
		return ErrResponseTooLarge
	}
	for i, ureq := range u.requests {
		if ureq.id == id {
			u.requests[i].resp = data[7 : respLen+7]
			return nil
		}
	}
	return ErrRequestNotFound
}

func (u *ProxyUser) responseWait(reqId uint16) ([]byte, error) {
	startTime := time.Now()
	for {
		for i, ureq := range u.requests {
			if ureq.id == reqId && ureq.resp != nil {
				defer u.sync.Unlock()
				u.sync.Lock()
				u.requests = append(u.requests[:i], u.requests[i+1:]...)
				return ureq.resp, nil
			}
		}
		if time.Since(startTime) > time.Millisecond*requestTimeout {
			return nil, ErrRequestTimeout
		}
	}
}
