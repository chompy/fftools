package main

import (
	"log"
	"runtime"
	"strconv"
	"strings"
)

// hexToInt converts hex string to an integer.
func hexToInt(hexString string) (int, error) {
	if hexString == "" {
		return 0, nil
	}
	output, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Println(err.Error(), hexString, fn, line)
	}
	return int(output), err
}

const roleTank = "tank"
const roleDps = "dps"
const roleHealer = "healer"

func jobGetRole(job string) string {
	switch strings.ToLower(job) {
	case "pld", "war", "drk", "gnb":
		{
			return roleTank
		}
	case "whm", "sch", "ast", "sge":
		{
			return roleHealer
		}
	}
	return roleDps
}
