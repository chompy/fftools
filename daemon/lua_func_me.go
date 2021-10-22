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
	L.Push(t)
	return 1
}

func findLocalPlayerLogLine(event *eventDispatch) {
	// TODO reset on new encounter?
	if localPlayerName != "" && localPlayerID > 0 {
		return
	}
	// collect player data from incomming log lines
	logLine := event.Data.(LogLine)
	parsedLogLine, err := ParseLogLine(logLine)
	if err != nil {
		logWarn(err.Error())
		return
	}
	switch parsedLogLine.Type {
	case LogTypeSingleTarget:
		{
			if localPlayerName != "" && parsedLogLine.AttackerName == localPlayerName {
				localPlayerID = parsedLogLine.AttackerID
				logInfo("Local player is %s (%d).", localPlayerName, localPlayerID)
			}
			break
		}
	case LogTypeChangePrimaryPlayer:
		{
			localPlayerName = parsedLogLine.AttackerName
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
