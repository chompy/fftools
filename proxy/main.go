package main

func main() {
	go webListen()
	if err := proxyListen(); err != nil {
		panic(err)
	}
}
