package main

import (
	"regexp"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

var regexPlayerCast = regexp.MustCompile(`You (use|cast) (.*).`)
var regexEnemyDamage = regexp.MustCompile(`(.*) takes ([0-9]*) damage.`)
var localPlayerID = 0
var localPlayerName = ""
var localPlayerCombatant = Combatant{ID: 0, Job: ""}
var localPlayerActions = make([][]interface{}, 0)
var playerActions = make([][]interface{}, 0)

func luaFuncMe(L *lua.LState) int {
	t := &lua.LTable{}
	t.RawSetString("name", lua.LString(localPlayerName))
	t.RawSetString("job", lua.LString(localPlayerCombatant.Job))
	t.RawSetString("id", lua.LNumber(localPlayerCombatant.ID))
	L.Push(t)
	return 1
}

func findLocalPlayerLogLine(event *eventDispatch) {
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
			playerActions = append(playerActions, []interface{}{
				parsedLogLine.AttackerName, parsedLogLine.AttackerID, parsedLogLine.AbilityName, parsedLogLine.Damage,
			})
			break
		}
	case LogTypeGameLog:
		{
			matches := regexPlayerCast.FindAllStringSubmatch(logLine.LogLine, -1)
			if len(matches) > 0 && len(matches[0]) > 0 {
				localPlayerActions = append(localPlayerActions, []interface{}{matches[0][2], ""})
			}
			matches = regexEnemyDamage.FindAllStringSubmatch(logLine.LogLine, -1)
			if len(matches) > 0 && len(matches[0]) > 0 && len(localPlayerActions) > 0 {
				dmg, err := strconv.Atoi(matches[0][2])
				if err != nil {
					logWarn(err.Error())
					return
				}
				localPlayerActions[len(localPlayerActions)-1][1] = dmg
			}
			break
		}
	}
	// see if there is enough info
	for _, localPlayerAction := range localPlayerActions {
		for _, playerAction := range playerActions {
			// compare skill name and damage delt to match player name
			if localPlayerAction[0] == playerAction[2] && localPlayerAction[1] == playerAction[3] {
				localPlayerName = playerAction[0].(string)
				localPlayerID = playerAction[1].(int)
				logInfo("Local player is %s (%d).", localPlayerName, localPlayerID)
				break
			}
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
