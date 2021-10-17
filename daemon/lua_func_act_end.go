package main

import (
	lua "github.com/yuin/gopher-lua"
)

func luaFuncActEnd(L *lua.LState) int {
	logLuaInfo(L, "ACT End encounter.")
	return 0
}

func init() {
	luaRegisterFunction("act_end", luaFuncActEnd)
}
