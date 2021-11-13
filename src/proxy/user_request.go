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

func (u *proxyUser) handleRequest(r *http.Request) (uint16, error) {
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
	log.Printf("[INFO] [%s] %s %s", u.uid, r.Method, ureq.url)
	u.requests = append(u.requests, ureq)
	out := make([]byte, 4)
	binary.LittleEndian.PutUint16(out[0:], ureq.id)
	binary.LittleEndian.PutUint16(out[2:], uint16(len(ureq.url)))
	out = append(out, []byte(ureq.url)...)
	if len(out) > requestMaxSize {
		return 0, ErrRequestTooLarge
	}
	if _, err := u.connection.Write(out); err != nil {
		return 0, err
	}
	return u.lastRequestId, nil
}

func (u *proxyUser) handleResponse(data []byte) error {
	id := binary.LittleEndian.Uint16(data[0:])
	respLen := binary.LittleEndian.Uint32(data[2:])
	if respLen+6 > responseMaxSize {
		return ErrResponseTooLarge
	}
	for i, ureq := range u.requests {
		if ureq.id == id {
			u.requests[i].resp = data[6 : respLen+6]
			return nil
		}
	}
	return ErrRequestNotFound
}

func (u *proxyUser) responseWait(reqId uint16) ([]byte, error) {
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
