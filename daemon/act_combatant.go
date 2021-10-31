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
