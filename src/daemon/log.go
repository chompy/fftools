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

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

const logPath = "data/daemon.log"

func logInfo(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func logWarn(msg string, args ...interface{}) {
	log.Printf("[WARN] "+msg, args...)
}

func logDebug(msg string, data interface{}) {
	log.Println("[DEBUG] "+msg, data)
}

func logPanic(err error) {
	log.Printf("[PANIC] " + err.Error())
	panic(err)
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
	scriptName := "?"
	if L != nil {
		scriptName = L.GetGlobal(luaGlobalScriptName).String()
	}
	msg = fmt.Sprintf("[%s] %s", scriptName, msg)
	logDebug(msg, data)
}
