package main

import (
	lua "github.com/yuin/gopher-lua"
)

func luaFuncLogInfo(L *lua.LState) int {
	msg := L.ToString(1)
	if msg == "" {
		return 0
	}
	logLuaInfo(L, msg)
	return 0
}

func luaFuncLogWarn(L *lua.LState) int {
	msg := L.ToString(1)
	if msg == "" {
		return 0
	}
	logLuaWarn(L, msg)
	return 0
}

func init() {
	luaRegisterFunction("log_info", luaFuncLogInfo)
	luaRegisterFunction("log_warn", luaFuncLogWarn)
}
