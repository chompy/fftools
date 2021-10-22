package main

import lua "github.com/yuin/gopher-lua"

func luaFuncEventDetach(L *lua.LState) int {
	ud := L.ToUserData(1)
	if ud != nil {
		listener := ud.Value.(*eventListener)
		eventListenerDetach(listener)
	}
	return 0
}

func luaEventDetachAllForState(L *lua.LState) {
	detachList := make([]*eventListener, 0)
	for _, eventListener := range eventListeners {
		switch eventListener.Data.(type) {
		case []interface{}:
			{
				data := eventListener.Data.([]interface{})
				if L == data[0].(*lua.LState) {
					detachList = append(detachList, eventListener)
				}
				break
			}
		}
	}
	for _, eventListener := range detachList {
		eventListenerDetach(eventListener)
	}
}

func init() {
	luaRegisterFunction("event_detach", luaFuncEventDetach)
}
