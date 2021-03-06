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
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func luaFuncActSay(L *lua.LState) int {
	text := L.ToString(1)
	if text == "" {
		return 0
	}
	logLuaInfo(L, "ACT TTS Say '%s.'", text)
	if err := actSay(text); err != nil {
		ls := L.GetGlobal(luaGlobalScriptData).(*lua.LUserData).Value.(*luaScript)
		ls.State = LuaScriptError
		ls.LastError = err
		logLuaWarn(L, err.Error())
		actError(err, ls.ScriptName)
	}
	return 0
}

func luaFuncActSayIf(L *lua.LState) int {
	// fetch condition
	cond := L.ToTable(2)
	name := ""
	job := ""
	role := ""
	id := 0
	if cond != nil && cond.Type() == lua.LTTable {
		nameL := cond.RawGetString("name")
		if nameL != nil && nameL.Type() == lua.LTString {
			name = string(nameL.(lua.LString))
		}
		jobL := cond.RawGetString("job")
		if jobL != nil && jobL.Type() == lua.LTString {
			job = string(jobL.(lua.LString))
		}
		roleL := cond.RawGetString("role")
		if roleL != nil && roleL.Type() == lua.LTString {
			role = string(roleL.(lua.LString))
		}
		idL := cond.RawGetString("id")
		if idL != nil && idL.Type() == lua.LTNumber {
			id = int(idL.(lua.LNumber))
		}
	}
	// no match
	if (name != "" && !strings.EqualFold(localPlayerName, name)) ||
		(job != "" && !strings.EqualFold(localPlayerCombatant.Job, job)) ||
		(role != "" && strings.EqualFold(jobGetRole(localPlayerCombatant.Job), role)) ||
		(id != 0 && localPlayerID != id) {
		return 0
	}
	if name == "" && role == "" && job == "" && id == 0 {
		return 0
	}
	return luaFuncActSay(L)
}

func init() {
	luaRegisterFunction("say", luaFuncActSay)
	luaRegisterFunction("say_if", luaFuncActSayIf)
}
