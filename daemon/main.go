package main

import "os"

func main() {
	luaEnableScripts()
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		testLogLines := testerParse(os.Stdin)
		if len(testLogLines) > 0 {
			testerReplay(testLogLines)
		}
	}
	go initWeb()
	actListenUDP()
}
