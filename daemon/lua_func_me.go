package main

import (
	lua "github.com/yuin/gopher-lua"
)

var localPlayerID = 0
var localPlayerName = ""
var localPlayerCombatant = Combatant{ID: 0, Job: ""}

func luaFuncMe(L *lua.LState) int {
	t := &lua.LTable{}
	t.RawSetString("name", lua.LString(localPlayerName))
	t.RawSetString("job", lua.LString(localPlayerCombatant.Job))
	t.RawSetString("id", lua.LNumber(localPlayerID))
	//t.RawSetString("name", lua.LString("Minda Silva"))
	//t.RawSetString("id", lua.LNumber(276036276))
	L.Push(t)
	return 1
}

func findLocalPlayerLogLine(event *eventDispatch) {
	// TODO reset on new encounter?
	if localPlayerName != "" && localPlayerID > 0 {
		return
	}
	// collect player data from log event
	logEvent := event.Data.(ParsedLogEvent)
	switch logEvent.Type {
	case LogTypeNetworkAbility:
		{
			if localPlayerName != "" && logEvent.Values["source_name"].(string) == localPlayerName {
				localPlayerID = logEvent.Values["source_id"].(int)
				logInfo("Local player is %s (%d).", localPlayerName, localPlayerID)
			}
			break
		}
	case LogTypeChangePrimaryPlayer:
		{
			localPlayerName = logEvent.Values["name"].(string)
			break
		}
	}
}

func findLocalPlayerCombatant(event *eventDispatch) {
	combatant := event.Data.(Combatant)
	if localPlayerID > 0 && combatant.ID == int32(localPlayerID) {
		localPlayerCombatant = combatant
	}
}

func init() {
	luaRegisterFunction("me", luaFuncMe)
	eventListenerAttach("act:log_line", findLocalPlayerLogLine)
	eventListenerAttach("act:combatant", findLocalPlayerCombatant)
}
