package main

import (
	lua "github.com/yuin/gopher-lua"
)

func luaFuncEventAttach(L *lua.LState) int {
	event := L.ToString(1)
	callback := L.ToFunction(2)
	if event == "" {
		return 0
	}
	logLuaInfo(L, "Attach event '%s.'", event)
	listener := eventListenerAttach(event, luaListenCallback)
	listener.Data = []interface{}{L, callback}
	L.Push(&lua.LUserData{Value: listener})
	return 1
}

func luaListenCallback(event *eventDispatch) {
	L := event.Listener.Data.([]interface{})[0].(*lua.LState)
	callback := event.Listener.Data.([]interface{})[1].(*lua.LFunction)
	logLuaDebug(L, "Recieved event.", event)
	argCount := 0
	var arg lua.LValue = nil
	switch event.Data.(type) {
	case LuaEncodable:
		{
			arg = event.Data.(LuaEncodable).ToLua()
			break
		}
	case string:
		{
			arg = lua.LString(event.Data.(string))
			break
		}
	}
	L.Push(callback)
	if arg != nil {
		argCount = 1
		L.Push(arg)
	}
	if err := L.PCall(argCount, 0, nil); err != nil {
		logLuaWarn(L, err.Error())
	}
}

func init() {
	luaRegisterFunction("event_attach", luaFuncEventAttach)
}
