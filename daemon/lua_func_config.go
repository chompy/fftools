package main

import (
	lua "github.com/yuin/gopher-lua"
)

func luaFuncConfigGet(L *lua.LState) int {
	key := L.ToString(1)
	if key == "" {
		return 0
	}
	luaScript := L.GetGlobal(luaGlobalScriptData).(*lua.LUserData).Value.(*luaScript)
	L.Push(valueGoToLua(luaScript.Config[key]))
	return 1
}

func init() {
	luaRegisterFunction("config_get", luaFuncConfigGet)
}
