package main

var currentEncounter Encounter
var currentZone = ""

func listenEncounter(event *eventDispatch) {
	if event.Data.(Encounter).ID != currentEncounter.ID {
		eventListenerDispatch("act:encounter:change", event.Data)
	}
	currentEncounter = event.Data.(Encounter)
	if currentZone != currentEncounter.Zone {
		currentZone = currentEncounter.Zone
		eventListenerDispatch("act:encounter:zone", currentZone)
	}
}

func init() {
	eventListenerAttach("act:encounter", listenEncounter)
}
