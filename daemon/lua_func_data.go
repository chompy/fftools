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
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

var luaData = make(map[string]map[string]interface{})

func dataGetPath(name string) string {
	return filepath.Join(getBasePath(), dataPath, name+".json")
}

func dataSaveFile(name string) error {
	jsonData, err := json.Marshal(luaData[name])
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(dataGetPath(name), jsonData, 0755); err != nil {
		return err
	}
	return nil
}

func dataLoadFile(name string) error {
	luaData[name] = make(map[string]interface{})
	jsonData, err := ioutil.ReadFile(dataGetPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			luaData[name] = make(map[string]interface{})
			return nil
		}
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}
	luaData[name] = data
	return nil
}

func luaFuncDataSet(L *lua.LState) int {
	key := L.ToString(1)
	if key == "" {
		logLuaWarn(L, "Invalid key for data_set.")
		return 0
	}
	val := L.Get(2)
	scriptName := L.GetGlobal(luaGlobalScriptName).String()
	if luaData[scriptName] == nil {
		if err := dataLoadFile(scriptName); err != nil {
			logLuaWarn(L, err.Error())
			return 0
		}
	}
	luaData[scriptName][key] = valueLuaToGo(val)
	if err := dataSaveFile(scriptName); err != nil {
		logLuaWarn(L, err.Error())
		return 0
	}
	L.Push(lua.LTrue)
	return 1
}

func luaFuncDataGet(L *lua.LState) int {
	scriptName := L.GetGlobal(luaGlobalScriptName).String()
	key := L.ToString(1)
	if key == "" {
		logLuaWarn(L, "Invalid key for data_get.")
		return 0
	}
	if luaData[scriptName] == nil {
		if err := dataLoadFile(scriptName); err != nil {
			logLuaWarn(L, err.Error())
			return 0
		}
	}
	L.Push(valueGoToLua(luaData[scriptName][key]))
	return 1
}

func init() {
	luaRegisterFunction("data_set", luaFuncDataSet)
	luaRegisterFunction("data_get", luaFuncDataGet)
}
