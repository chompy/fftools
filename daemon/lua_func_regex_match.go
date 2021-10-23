package main

import (
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

var luaRegexList map[string]*regexp.Regexp = nil

func luaFuncRegexMatch(L *lua.LState) int {

	if luaRegexList == nil {
		luaRegexList = make(map[string]*regexp.Regexp)
	}

	regexStr := L.ToString(1)
	if regexStr == "" {
		return 0
	}
	if luaRegexList[regexStr] == nil {
		var err error
		luaRegexList[regexStr], err = regexp.Compile(regexStr)
		if err != nil {
			logLuaWarn(L, err.Error())
			scriptName := L.GetGlobal(luaGlobalScriptName).String()
			actError(err, scriptName)
			return 0
		}
	}
	matchStr := L.ToString(2)
	if matchStr == "" {
		return 0
	}
	matches := luaRegexList[regexStr].FindAllStringSubmatch(matchStr, -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		return 0
	}
	t := &lua.LTable{}
	for _, match := range matches[0] {
		t.Append(lua.LString(match))
	}
	L.Push(t)
	return 1
}

func init() {
	luaRegisterFunction("regex_match", luaFuncRegexMatch)
}
