package main

import (
	"net"
	"time"
)

const proxyUserTTL = 86400

type proxyUser struct {
	uid             string
	connection      net.Conn
	lastRequestTime time.Time
}

var proxyUsers = make([]*proxyUser, 0)

func addProxyUser(conn net.Conn) *proxyUser {
	uid := uidFromConn(conn)
	u := getProxyUser(uid)
	if u != nil {
		u.connection = conn
		return u
	}
	u = &proxyUser{
		uid:             uid,
		connection:      conn,
		lastRequestTime: time.Now(),
	}
	proxyUsers = append(proxyUsers, u)
	// clean up routine
	go func(index int) {
		for {
			time.Sleep(time.Hour)
			if time.Since(u.lastRequestTime) > time.Second*proxyUserTTL {
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
