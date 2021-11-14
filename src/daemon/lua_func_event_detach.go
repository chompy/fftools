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
