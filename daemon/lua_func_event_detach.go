package main

import lua "github.com/yuin/gopher-lua"

func luaFuncEventDetach(L *lua.LState) int {
	ud := L.ToUserData(1)
	if ud != nil {
		listener := ud.Value.(*eventListener)
		eventListenerDetach(listener)
	}
	return 0
}

func init() {
	luaRegisterFunction("event_detach", luaFuncEventDetach)
}
