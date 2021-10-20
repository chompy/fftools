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

func init() {
	luaRegisterFunction("act_encounter", luaFuncActEncounter)
	luaRegisterFunction("act_encounter_time", luaFuncActEncounterTime)
}
