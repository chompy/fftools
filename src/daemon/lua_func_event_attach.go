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
	ls := L.GetGlobal(luaGlobalScriptData).(*lua.LUserData).Value.(*luaScript)
	ls.Lock.Lock()
	defer ls.Lock.Unlock()
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
	case lua.LValue:
		{
			arg = event.Data.(lua.LValue)
			break
		}
	}
	L.SetTop(0)
	L.Push(callback)
	if arg != nil {
		argCount = 1
		L.Push(arg)
	}
	if err := L.PCall(argCount, 0, nil); err != nil {
		ls.State = LuaScriptError
		ls.LastError = err
		logLuaWarn(L, err.Error())
		actError(err, ls.ScriptName)
	}
}

func init() {
	luaRegisterFunction("event_attach", luaFuncEventAttach)
}
