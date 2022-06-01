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
			ls := L.GetGlobal(luaGlobalScriptData).(*lua.LUserData).Value.(*luaScript)
			ls.State = LuaScriptError
			ls.LastError = err
			logLuaWarn(L, err.Error())
			actError(err, ls.ScriptName)
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
