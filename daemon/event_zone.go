package main

var currentZone = ""

func listenZoneChange(event *eventDispatch) {
	encounter := event.Data.(Encounter)
	if currentZone != encounter.Zone {
		currentZone = encounter.Zone
		eventListenerDispatch("zone", currentZone)
	}
}

func init() {
	eventListenerAttach("act:encounter", listenZoneChange)
}
