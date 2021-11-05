package main

import (
	"os"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

const luaGlobalScriptName = "_script_name"
const luaGlobalScriptData = "_script_data"

type luaScript struct {
	Name       string
	ScriptName string
	Desc       string
	Path       string
	LastError  error
	Enabled    bool
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
		Enabled:    false,
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
	ls.close()
	ls.LastError = nil
	// config
	var err error
	ls.Config, err = configLoadScriptConfig(ls.ScriptName)
	if err != nil && !os.IsNotExist(err) {
		logLuaWarn(ls.L, err.Error())
	}
	// init lua
	ls.L = lua.NewState()
	// set global script name var
	ls.L.SetGlobal(luaGlobalScriptName, lua.LString(ls.ScriptName))
	ls.L.SetGlobal(luaGlobalScriptData, &lua.LUserData{Value: ls})
	// load script
	if err := ls.L.DoFile(ls.Path); err != nil {
		ls.LastError = err
		return err
	}
	// set global functions
	funcTable := &lua.LTable{}
	for name, function := range luaFuncs {
		funcTable.RawSetString(name, ls.L.NewFunction(function))
		ls.L.SetGlobal(name, ls.L.NewFunction(function))
	}
	ls.L.SetGlobal("ffl", funcTable)

	logLuaInfo(ls.L, "Loaded.")
	return nil
}

func (ls *luaScript) close() {
	if ls.L != nil {
		logLuaInfo(ls.L, "Unloaded.")
		luaEventDetachAllForState(ls.L)
		ls.L.Close()
		ls.L = nil
	}
}

func (ls *luaScript) init() error {
	if !ls.Enabled {
		return nil
	}
	ls.Lock.Lock()
	defer ls.Lock.Unlock()
	logLuaInfo(ls.L, "Enabled.")
	initFunc := ls.L.GetGlobal("init")
	ls.L.SetTop(0)
	ls.L.Push(initFunc)
	err := ls.L.PCall(0, 0, nil)
	if err != nil {
		ls.LastError = err
		ls.Enabled = false
		logLuaWarn(ls.L, err.Error())
		actError(err, ls.ScriptName)
		return err
	}
	return nil
}

func (ls *luaScript) info() error {
	ls.Lock.Lock()
	defer ls.Lock.Unlock()
	ls.L.SetTop(0)
	infoFunc := ls.L.GetGlobal("info")
	ls.L.Push(infoFunc)
	if err := ls.L.PCall(0, 1, nil); err != nil {
		ls.LastError = err
		logLuaWarn(ls.L, err.Error())
		return err
	}
	infoTable := ls.L.ToTable(1)
	name := infoTable.RawGetString("name")
	desc := infoTable.RawGetString("desc")
	ls.Name = string(name.(lua.LString))
	ls.Desc = string(desc.(lua.LString))
	return nil
}
