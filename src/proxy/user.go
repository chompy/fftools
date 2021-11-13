package main

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const proxyUserTTL = 86400

type proxyUser struct {
	uid             string
	connection      net.Conn
	lastRequestTime time.Time
	requests        []userRequest
	lastRequestId   uint16
	sync            sync.Mutex
}

var proxyUsers = make([]*proxyUser, 0)

func addProxyUser(uid string, conn net.Conn) *proxyUser {
	u := getProxyUser(uid)
	if u != nil {
		u.connection = conn
		return u
	}
	u = &proxyUser{
		uid:             uid,
		connection:      conn,
		lastRequestTime: time.Now(),
		requests:        make([]userRequest, 0),
	}
	proxyUsers = append(proxyUsers, u)
	go func(index int) {
		buf := make([]byte, responseMaxSize)
		for {

			// check response
			u := proxyUsers[index]

			n, err := u.connection.Read(buf)
			if err != nil {
				log.Printf("[WARN] proxy read response :: %s", err.Error())
				if errors.Is(err, io.EOF) {
					u.lastRequestTime = time.Time{}
				}
			}
			if n > 0 {
				if err := u.handleResponse(buf); err != nil {
					log.Printf("[WARN] proxy handle response :: %s", err.Error())
				}
			}
			// clean up
			if time.Since(u.lastRequestTime) > time.Second*proxyUserTTL {
				log.Printf("[INFO] Clean up '%s.'", u.uid)
				u.connection.Close()
				proxyUsers = append(proxyUsers[index:], proxyUsers[:index+1]...)
				return
			}
		}
	}(len(proxyUsers) - 1)

	return u
}

func getProxyUser(uid string) *proxyUser {
	for _, u := range proxyUsers {
		if u.uid == uid {
			return u
		}
	}
	return nil
}
