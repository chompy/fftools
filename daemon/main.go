package main

func main() {

	// init
	eventListenerReset()

	for _, name := range luaGetAvailableScripts() {
		ls, err := luaLoadScript(name)
		if err != nil {
			panic(err)
		}
		ls.init()
	}

	ListenUDP()
}
