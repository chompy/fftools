/*
This file is part of FFTools.

FFTools is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

FFTools is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with FFTools.  If not, see <https://www.gnu.org/licenses/>.
*/

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
