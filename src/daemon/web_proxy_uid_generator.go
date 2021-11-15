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
	rand.Seed(time.Now().UnixMicro())
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
