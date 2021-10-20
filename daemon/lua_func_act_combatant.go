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
	luaRegisterFunction("act_combatants", luaFuncActCombatants)
	luaRegisterFunction("act_combatant_from_id", luaFuncActCombatantFromId)
	luaRegisterFunction("act_combatant_from_name", luaFuncActCombatantFromName)
}
