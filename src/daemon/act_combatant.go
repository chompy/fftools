/*
This file is part of FF Lua.

FF Lua is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

FF Lua is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with FF Lua.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

var currentCombatants = make([]Combatant, 0)
var combatantNameLookup = make(map[int]string)

func listenCombatant(event *eventDispatch) {
	combatant := event.Data.(Combatant)
	if combatant.ID > 999999999 ||
		combatant.Job == "" ||
		len(combatant.Job) > 3 ||
		combatantNameLookup[int(combatant.ID)] == "" {
		return
	}
	for i := range currentCombatants {
		if currentCombatants[i].ID == combatant.ID {
			currentCombatants[i] = combatant
			currentCombatants[i].Name = combatantNameLookup[int(combatant.ID)]
			return
		}
	}
	currentCombatants = append(currentCombatants, combatant)
	currentCombatants[len(currentCombatants)-1].Name = combatantNameLookup[int(combatant.ID)]
	logInfo("Register combatant %s (#%d).", combatantNameLookup[int(combatant.ID)], combatant.ID)
}

func listenCombatantLog(event *eventDispatch) {
	log := event.Data.(ParsedLogEvent)
	if log.Type != LogTypeNetworkAbility || log.Values["source_name"].(string) == "" {
		return
	}
	if combatantNameLookup[log.Values["source_id"].(int)] == "" {
		combatantNameLookup[log.Values["source_id"].(int)] = log.Values["source_name"].(string)
	}
}

func listenEncounterChangeResetCombatants(event *eventDispatch) {
	currentCombatants = make([]Combatant, 0)
	combatantNameLookup = make(map[int]string)
}

func init() {
	eventListenerAttach("act:log_line", listenCombatantLog)
	eventListenerAttach("act:combatant", listenCombatant)
	eventListenerAttach("act:encounter:change", listenEncounterChangeResetCombatants)
}
