package main

import (
	lua "github.com/yuin/gopher-lua"
)

func luaFuncXivParseLogLine(L *lua.LState) int {
	// convert lua value back to log line + ensure pased lua arg is log line
	logLineTable := L.ToTable(1)
	if logLineTable == nil {
		logLuaWarn(L, "Value passed to 'xiv_parse_log_line' is not a log line.")
		return 0
	}
	logLineL, ok := logLineTable.RawGetString("__goobject").(*lua.LUserData)
	if !ok {
		logLuaWarn(L, "Value passed to 'xiv_parse_log_line' is not a log line.")
		return 0
	}
	logLine, ok := logLineL.Value.(LogLine)
	if !ok {
		logLuaWarn(L, "Value passed to 'xiv_parse_log_line' is not a log line.")
		return 0
	}
	// parse the log line
	parsedLogLine, err := ParseLogLine(logLine)
	if err != nil {
		logLuaWarn(L, err.Error())
		return 0
	}
	// convert parsed log line to lua
	L.Push(parsedLogLine.ToLua())
	return 1
}

func init() {
	luaRegisterFunction("xiv_parse_log_line", luaFuncXivParseLogLine)
}
