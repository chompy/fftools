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
	log.Println("[DEBUG] "+msg, data)
}

func logLuaInfo(L *lua.LState, msg string, args ...interface{}) {
	scriptName := L.GetGlobal(luaGlobalScriptName).String()
	msg = fmt.Sprintf("[%s] %s", scriptName, msg)
	logInfo(msg, args...)
}

func logLuaWarn(L *lua.LState, msg string, args ...interface{}) {
	scriptName := L.GetGlobal(luaGlobalScriptName).String()
	msg = fmt.Sprintf("[%s] %s", scriptName, msg)
	logWarn(msg, args...)
}

func logLuaDebug(L *lua.LState, msg string, data interface{}) {
	scriptName := L.GetGlobal(luaGlobalScriptName).String()
	msg = fmt.Sprintf("[%s] %s", scriptName, msg)
	logDebug(msg, data)
}
