package main

import (
	"bufio"
	"log"
	"net"
)

const webProxyAddress = "localhost:31595"

func webProxyConnect() error {

	conn, err := net.Dial("tcp", webProxyAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	buf := bufio.NewReader(conn)
	for {
		b, _ := buf.ReadByte()
		log.Println(int(b))
	}

}
