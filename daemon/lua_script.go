package main

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

type luaScript struct {
	Name       string
	ScriptName string
	Desc       string
	Path       string
	LastError  error
	Enabled    bool
	Config     map[string]interface{}
	L          *lua.LState
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
	ls.Config, err = configLoadScriptConfig(name)
	if err != nil && !os.IsNotExist(err) {
		logLuaWarn(ls.L, err.Error())
	}
	if err := ls.info(); err != nil {
		return ls, err
	}
	return ls, nil
}

func (ls *luaScript) load() error {
	ls.close()
	ls.LastError = nil
	// init lua, load script file
	ls.L = lua.NewState()
	if err := ls.L.DoFile(ls.Path); err != nil {
		ls.close()
		ls.LastError = err
		return err
	}
	// set globals variables
	ls.L.SetGlobal(luaGlobalScriptName, lua.LString(ls.ScriptName))
	// set global functions
	for name, function := range luaFuncs {
		ls.L.SetGlobal(name, ls.L.NewFunction(function))
	}
	logLuaInfo(ls.L, "Loaded.")
	return nil
}

func (ls *luaScript) close() {
	if ls.L != nil {
		logLuaInfo(ls.L, "Disabled.")
		luaEventDetachAllForState(ls.L)
		ls.L.Close()
		ls.L = nil
	}
}

func (ls *luaScript) init() error {
	logLuaInfo(ls.L, "Enabled.")
	initFunc := ls.L.GetGlobal("init")
	ls.L.Push(initFunc)
	err := ls.L.PCall(0, 0, nil)
	if err != nil {
		ls.LastError = err
		actError(err, ls.ScriptName)
		return err
	}
	return nil
}

func (ls *luaScript) info() error {
	infoFunc := ls.L.GetGlobal("info")
	ls.L.Push(infoFunc)
	if err := ls.L.PCall(0, 1, nil); err != nil {
		ls.LastError = err
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
