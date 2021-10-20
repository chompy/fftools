package main

var currentCombatants []Combatant = make([]Combatant, 0)

func listenCombatant(event *eventDispatch) {
	combatant := event.Data.(Combatant)
	for i := range currentCombatants {
		if currentCombatants[i].ID == combatant.ID {
			currentCombatants[i] = combatant
			return
		}
	}
	currentCombatants = append(currentCombatants, combatant)
}

func listenEncounterChangeResetCombatants(event *eventDispatch) {
	currentCombatants = make([]Combatant, 0)
}

func init() {
	eventListenerAttach("act:combatant", listenCombatant)
	eventListenerAttach("act:encounter:change", listenEncounterChangeResetCombatants)
}
