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
	lua "github.com/yuin/gopher-lua"
)

func luaFuncActCombatants(L *lua.LState) int {
	ltable := L.NewTable()
	for _, combatant := range currentCombatants {
		ltable.Append(combatant.ToLua())
	}
	L.Push(ltable)
	return 1
}

func luaFuncActCombatantFromId(L *lua.LState) int {
	id := L.ToInt(1)
	for _, combatant := range currentCombatants {
		if combatant.ID == int32(id) {
			L.Push(combatant.ToLua())
			return 1
		}
	}
	return 0
}

func luaFuncActCombatantFromName(L *lua.LState) int {
	name := L.ToString(1)
	for _, combatant := range currentCombatants {
		if combatant.Name == name {
			L.Push(combatant.ToLua())
			return 1
		}
	}
	return 0
}

func init() {
	luaRegisterFunction("combatants", luaFuncActCombatants)
	luaRegisterFunction("combatant_from_id", luaFuncActCombatantFromId)
	luaRegisterFunction("combatant_from_name", luaFuncActCombatantFromName)
}
