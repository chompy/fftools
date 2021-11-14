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
	"errors"
	"net/http"
	"os"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

const luaGlobalScriptName = "_script_name"
const luaGlobalScriptData = "_script_data"

const (
	LuaScriptInactive = 0
	LuaScriptActive   = 1
	LuaScriptError    = 2
)

type luaScript struct {
	Name       string
	ScriptName string
	Desc       string
	Path       string
	LastError  error
	State      int
	Config     map[string]interface{}
	L          *lua.LState
	Lock       sync.Mutex
}

func luaLoadScript(name string) (*luaScript, error) {
	pathTo, err := luaGetPathToScript(name)
	if err != nil {
		return nil, err
	}
	ls := &luaScript{
		Name:       name,
		Desc:       "N/A",
		ScriptName: name,
		Path:       pathTo,
		State:      LuaScriptInactive,
	}
	if err := ls.load(); err != nil {
		return ls, err
	}
	if err := ls.info(); err != nil {
		return ls, err
	}
	return ls, nil
}

func (ls *luaScript) load() error {
	ls.Lock.Lock()
	defer ls.Lock.Unlock()
	ls.unload()
	ls.LastError = nil
	// config
	var err error
	ls.Config, err = configLoadScriptConfig(ls.ScriptName)
	if err != nil && !os.IsNotExist(err) && !errors.Is(err, ErrDefaultConfigNotFound) {
		ls.LastError = err
		ls.State = LuaScriptError
		logLuaWarn(ls.L, err.Error())
		actError(err, ls.ScriptName)
		return err
	}
	// init lua
	ls.L = lua.NewState()
	// set global script name var
	ls.L.SetGlobal(luaGlobalScriptName, lua.LString(ls.ScriptName))
	ls.L.SetGlobal(luaGlobalScriptData, &lua.LUserData{Value: ls})
	logLuaInfo(ls.L, "Load.")
	// load script
	if err := ls.L.DoFile(ls.Path); err != nil {
		ls.LastError = err
		ls.State = LuaScriptError
		logLuaWarn(ls.L, err.Error())
		actError(err, ls.ScriptName)
		return err
	}
	// set global functions
	funcTable := &lua.LTable{}
	for name, function := range luaFuncs {
		funcTable.RawSetString(name, ls.L.NewFunction(function))
		ls.L.SetGlobal("fft_"+name, ls.L.NewFunction(function))
	}
	ls.L.SetGlobal("fft", funcTable)
	return nil
}

func (ls *luaScript) reload() error {
	ls.unload()
	if err := ls.load(); err != nil {
		ls.LastError = err
		ls.State = LuaScriptError
		logLuaWarn(ls.L, err.Error())
		actError(err, ls.ScriptName)
		return err
	}
	if err := ls.info(); err != nil {
		ls.LastError = err
		ls.State = LuaScriptError
		logLuaWarn(ls.L, err.Error())
		actError(err, ls.ScriptName)
		return err
	}
	return nil
}

func (ls *luaScript) unload() {
	if ls.L != nil {
		logLuaInfo(ls.L, "Deactivate.")
		luaEventDetachAllForState(ls.L)
		ls.L.Close()
		ls.L = nil
		ls.State = LuaScriptInactive
		ls.LastError = nil
	}
}

func (ls *luaScript) activate() error {
	if ls.L == nil {
		return ErrLuaScriptNotLoaded
	}
	if ls.State == LuaScriptError {
		return ErrLuaScriptInError
	}
	ls.Lock.Lock()
	defer ls.Lock.Unlock()
	logLuaInfo(ls.L, "Activate.")
	initFunc := ls.L.GetGlobal("init")
	ls.L.SetTop(0)
	ls.L.Push(initFunc)
	err := ls.L.PCall(0, 0, nil)
	if err != nil {
		ls.LastError = err
		ls.State = LuaScriptError
		logLuaWarn(ls.L, err.Error())
		actError(err, ls.ScriptName)
		return err
	}
	ls.State = LuaScriptActive
	return nil
}

func (ls *luaScript) info() error {
	if ls.L == nil {
		return ErrLuaScriptNotLoaded
	}
	if ls.State == LuaScriptError {
		return ErrLuaScriptInError
	}
	ls.Lock.Lock()
	defer ls.Lock.Unlock()
	ls.L.SetTop(0)
	infoFunc := ls.L.GetGlobal("info")
	ls.L.Push(infoFunc)
	if err := ls.L.PCall(0, 1, nil); err != nil {
		ls.LastError = err
		ls.State = LuaScriptError
		logLuaWarn(ls.L, err.Error())
		actError(err, ls.ScriptName)
		return err
	}
	infoTable := ls.L.ToTable(1)
	name := infoTable.RawGetString("name")
	desc := infoTable.RawGetString("desc")
	ls.Name = string(name.(lua.LString))
	ls.Desc = string(desc.(lua.LString))
	return nil
}

// Web calls the lua script's web function if it exists.
func (ls *luaScript) Web(r *http.Request) (interface{}, error) {
	ls.L.SetTop(0)
	// push request values + query params
	luaRequest := &lua.LTable{}
	luaRequest.RawSetString("url", lua.LString(r.URL.String()))
	luaRequest.RawSetString("host", lua.LString(r.URL.Host))
	luaRequest.RawSetString("hostname", lua.LString(r.URL.Hostname()))
	luaRequest.RawSetString("port", lua.LString(r.URL.Port()))
	luaRequest.RawSetString("path", lua.LString(r.URL.Path))
	queryTable := &lua.LTable{}
	for k, v := range r.URL.Query() {
		queryTable.RawSetString(k, valueGoToLua(v))
	}
	luaRequest.RawSetString("query", queryTable)

	luaFunc := ls.L.GetGlobal("web")
	if luaFunc.Type() != lua.LTFunction {
		return nil, nil
	}
	ls.L.Push(luaFunc)
	ls.L.Push(luaRequest)
	if err := ls.L.PCall(1, 1, nil); err != nil {
		ls.LastError = err
		ls.State = LuaScriptError
		logLuaWarn(ls.L, err.Error())
		actError(err, ls.ScriptName)
		return nil, err
	}
	return valueLuaToGo(ls.L.Get(1)), nil
}
