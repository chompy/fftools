package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var proxyUidRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const proxyUidLen = 8
const secretSalt = "U2?6M?AEq+XZ--!232Dvckcla/.,;sDfd"

func proxyRandUid(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = proxyUidRunes[rand.Intn(len(proxyUidRunes))]
	}
	return string(b)
}

func webProxyGenerateCreds() (string, string) {
	uid := proxyRandUid(proxyUidLen)
	secretGenStr := fmt.Sprintf("%s-%d-%d-%s", secretSalt, time.Now().UnixMicro(), rand.Uint64(), uid)
	secret := fmt.Sprintf("%x", sha256.Sum256([]byte(secretGenStr)))
	return uid, secret
}

func webProxySendCreds(uid string, secret string, conn net.Conn) error {
	data := make([]byte, 1)
	data[0] = proxyMsgLogin
	data = append(data, []byte(uid)...)
	data = append(data, []byte(secret)...)
	_, err := conn.Write(data)
	return err
}
