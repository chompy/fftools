package main

func main() {
	luaEnableScripts()
	go initWeb()
	actListenUDP()
}
