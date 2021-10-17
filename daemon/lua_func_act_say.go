package main

import (
	lua "github.com/yuin/gopher-lua"
)

func luaFuncActSay(L *lua.LState) int {
	arg1Val := L.ToString(1)
	if arg1Val == "" {
		return 0
	}
	logLuaInfo(L, "ACT TTS Say '%s.'", arg1Val)
	return 0
}

func init() {
	luaRegisterFunction("act_say", luaFuncActSay)
}
