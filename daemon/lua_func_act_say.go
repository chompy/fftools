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
		logLuaWarn(L, err.Error())
		scriptName := L.GetGlobal(luaGlobalScriptName).String()
		actError(err, scriptName)
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
	if cond != nil {
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
	return luaFuncActSayIf(L)
}

func init() {
	luaRegisterFunction("act_say", luaFuncActSay)
	luaRegisterFunction("act_say_if", luaFuncActSayIf)
}
