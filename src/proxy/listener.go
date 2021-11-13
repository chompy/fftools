package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/martinlindhe/base36"
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
		uid := uidFromConn(conn)
		log.Printf("[INFO] Connection from %s (%s).", conn.RemoteAddr().String(), uid)
		if err != nil {
			log.Printf("[WARN] proxyListen :: %s", err.Error())
			continue
		}
		addProxyUser(uid, conn)
	}
}

func uidFromConn(conn net.Conn) string {
	addr := strings.Split(conn.RemoteAddr().String(), ":")
	ipStr := strings.Trim(strings.Join(addr[:len(addr)-1], ":"), "[]")
	ip := net.ParseIP(ipStr)
	/*port, err := strconv.Atoi(addr[len(addr)-1])
	if err != nil {
		return ""
	}*/
	ipHash := sha256.Sum256(ip)
	ipHashBytes := make([]byte, 32)
	for i := range ipHash {
		ipHashBytes[i] = ipHash[i]
	}
	return strings.ToLower(base36.EncodeBytes(ipHashBytes))[0:8]
}
