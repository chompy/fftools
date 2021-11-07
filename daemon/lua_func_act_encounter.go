package main

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

func luaFuncActEncounter(L *lua.LState) int {
	L.Push(currentEncounter.ToLua())
	return 1
}

func luaFuncActEncounterTime(L *lua.LState) int {
	dur := currentEncounter.EndTime.Sub(currentEncounter.StartTime)
	if currentEncounter.Active {
		dur = time.Since(currentEncounter.StartTime)
	}
	L.Push(lua.LNumber(dur.Milliseconds()))
	return 1
}

func luaFuncActEnd(L *lua.LState) int {
	logLuaInfo(L, "ACT End encounter.")
	if err := actEnd(); err != nil {
		ls := L.GetGlobal(luaGlobalScriptData).(*lua.LUserData).Value.(*luaScript)
		ls.State = LuaScriptError
		ls.LastError = err
		logLuaWarn(L, err.Error())
		actError(err, ls.ScriptName)
	}
	return 0
}

func init() {
	luaRegisterFunction("encounter", luaFuncActEncounter)
	luaRegisterFunction("encounter_time", luaFuncActEncounterTime)
	luaRegisterFunction("encounter_end", luaFuncActEnd)
}
