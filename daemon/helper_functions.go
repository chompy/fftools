package main

import (
	"log"
	"runtime"
	"strconv"
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
