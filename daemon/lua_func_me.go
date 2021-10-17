package main

import (
	lua "github.com/yuin/gopher-lua"
)

func luaFuncMe(L *lua.LState) int {
	t := &lua.LTable{}
	// TODO
	t.RawSetString("name", lua.LString("Qunara Sivra"))
	L.Push(t)
	return 1
}

func init() {
	luaRegisterFunction("me", luaFuncMe)
}
