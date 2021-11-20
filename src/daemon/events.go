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
	// don't log the log_line event as it's too noisy
	if event != "act:log_line" {
		logDebug("Dispatch event.", eventDispatch{Event: event, Data: data})
	}
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
