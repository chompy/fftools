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

func luaFuncConfigGet(L *lua.LState) int {
	key := L.ToString(1)
	if key == "" {
		return 0
	}
	luaScript := L.GetGlobal(luaGlobalScriptData).(*lua.LUserData).Value.(*luaScript)
	L.Push(valueGoToLua(luaScript.Config[key]))
	return 1
}

func init() {
	luaRegisterFunction("config_get", luaFuncConfigGet)
}
