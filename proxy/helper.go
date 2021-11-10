package main

import (
	"net"
	"strconv"
	"strings"

	"github.com/martinlindhe/base36"
)

func uidFromConn(conn net.Conn) string {
	addr := strings.Split(conn.RemoteAddr().String(), ":")
	ipStr := strings.Trim(strings.Join(addr[:len(addr)-1], ":"), "[]")
	ip := net.ParseIP(ipStr)
	port, err := strconv.Atoi(addr[len(addr)-1])
	if err != nil {
		return ""
	}
	return strings.ToLower(base36.EncodeBytes(ip) + "-" + base36.Encode(uint64(port)))
}

func addressFromUid(uid string) (net.IP, uint16) {
	s := strings.Split(strings.ToUpper(uid), "-")
	if len(s) != 2 {
		return net.IP{}, 0
	}
	ip := net.IP(base36.DecodeToBytes(s[0]))
	port := uint16(base36.Decode(s[1]))
	return ip, port
}
