package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

const proxyPort = 31595

func proxyListen() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", proxyPort))
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		log.Printf("[INFO] Connection from %s.", conn.RemoteAddr().String())
		if err != nil {
			log.Printf("[WARN] %s", err.Error())
			continue
		}
		addProxyUser(conn)
	}
}

func proxyRequest(user *proxyUser, r *http.Request) error {

	rawReq, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}

	user.connection.Write(rawReq)
	user.connection.R

}

func proxyResponse() {

}d
