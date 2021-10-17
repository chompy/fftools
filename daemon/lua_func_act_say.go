package main

import (
	lua "github.com/yuin/gopher-lua"
)

func luaFuncActSay(L *lua.LState) int {
	text := L.ToString(1)
	if text == "" {
		return 0
	}
	logLuaInfo(L, "ACT TTS Say '%s.'", text)
	if err := actSay(text); err != nil {
		logLuaWarn(L, err.Error())
	}
	return 0
}

func init() {
	luaRegisterFunction("act_say", luaFuncActSay)
}
