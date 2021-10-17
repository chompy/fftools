package main

func main() {

	// init
	//eventListenerReset()

	for _, name := range luaGetEnabledScripts() {
		ls, err := luaLoadScript(name)
		if err != nil {
			panic(err)
		}
		ls.init()
	}

	actListenUDP()
}
