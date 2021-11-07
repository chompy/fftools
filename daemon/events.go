package main

type eventDispatch struct {
	Event    string
	Data     interface{}
	Listener *eventListener
}

type eventListener struct {
	Event    string
	Callback func(*eventDispatch)
	Data     interface{}
}

var eventListeners = make([]*eventListener, 0)

/*func eventListenerReset() {
	eventListeners = make([]*eventListener, 0)
}*/

func eventListenerAttach(event string, callback func(*eventDispatch)) *eventListener {
	if event == "" {
		return nil
	}
	listener := &eventListener{
		Event:    event,
		Callback: callback,
	}
	eventListeners = append(eventListeners, listener)
	return listener
}

func eventListenerDetach(listener *eventListener) {
	for i := range eventListeners {
		if listener == eventListeners[i] {
			logDebug("Detach event listener.", listener)
			eventListeners = append(eventListeners[:i], eventListeners[i+1:]...)
			return
		}
	}
}

func eventListenerDispatch(event string, data interface{}) {
	if event == "" {
		return
	}
	logDebug("Dispatch event.", eventDispatch{Event: event, Data: data})
	for _, listener := range eventListeners {
		if listener.Event == event {
			ed := &eventDispatch{
				Event:    event,
				Data:     data,
				Listener: listener,
			}
			listener.Callback(ed)
		}
	}
}
