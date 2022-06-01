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
	// collect player data from log event
	logEvent := event.Data.(ParsedLogEvent)
	switch logEvent.Type {
	case LogTypeNetworkUpdateHP:
		{
			if localPlayerName != "" && logEvent.Values["target_name"].(string) == localPlayerName {
				newLocalPlayerID := logEvent.Values["target_id"].(int)
				if localPlayerID != newLocalPlayerID {
					logInfo("Set active player to %s (#%d).", localPlayerName, newLocalPlayerID)
				}
				localPlayerID = newLocalPlayerID
			}
			break
		}
	case LogTypeChangePrimaryPlayer:
		{
			localPlayerName = logEvent.Values["target_name"].(string)
			// reset
			localPlayerID = 0
			localPlayerCombatant = Combatant{ID: 0, Job: ""}
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
