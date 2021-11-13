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
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

func logInfo(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func logWarn(msg string, args ...interface{}) {
	log.Printf("[WARN] "+msg, args...)
}

func logDebug(msg string, data interface{}) {
	//log.Println("[DEBUG] "+msg, data)
}

func logLuaInfo(L *lua.LState, msg string, args ...interface{}) {
	scriptName := "?"
	if L != nil {
		scriptName = L.GetGlobal(luaGlobalScriptName).String()
	}
	msg = fmt.Sprintf("[%s] %s", scriptName, msg)
	logInfo(msg, args...)
}

func logLuaWarn(L *lua.LState, msg string, args ...interface{}) {
	scriptName := "?"
	if L != nil {
		scriptName = L.GetGlobal(luaGlobalScriptName).String()
	}
	msg = fmt.Sprintf("[%s] %s", scriptName, msg)
	logWarn(msg, args...)
}

func logLuaDebug(L *lua.LState, msg string, data interface{}) {
	//logLua(L, "[DEBUG] "+msg+" %s", data)
}
