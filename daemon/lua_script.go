package main

import (
	lua "github.com/yuin/gopher-lua"
)

type luaScript struct {
	Name       string
	ScriptName string
	Desc       string
	Path       string
	L          *lua.LState
}

func luaLoadScript(name string) (*luaScript, error) {
	pathTo, err := luaGetPathToScript(name)
	if err != nil {
		return nil, err
	}
	L := lua.NewState()
	err = L.DoFile(pathTo)
	if err != nil {
		L.Close()
		return nil, err
	}
	// set globals variables
	L.SetGlobal(luaGlobalScriptName, lua.LString(name))
	// set global functions
	for name, function := range luaFuncs {
		L.SetGlobal(name, L.NewFunction(function))
	}
	ls := &luaScript{
		ScriptName: name,
		Path:       pathTo,
		L:          L,
	}
	if err := ls.info(); err != nil {
		return nil, err
	}
	return ls, nil
}

func (ls *luaScript) close() {
	luaEventDetachAllForState(ls.L)
	ls.L.Close()
}

func (ls *luaScript) init() error {
	logInfo("Init '%s' Lua script.", ls.ScriptName)
	initFunc := ls.L.GetGlobal("init")
	ls.L.Push(initFunc)
	return ls.L.PCall(0, 0, nil)
}

func (ls *luaScript) info() error {
	infoFunc := ls.L.GetGlobal("info")
	ls.L.Push(infoFunc)
	if err := ls.L.PCall(0, 1, nil); err != nil {
		return err
	}
	infoTable := ls.L.ToTable(1)
	name := infoTable.RawGetString("name")
	desc := infoTable.RawGetString("desc")
	ls.Name = string(name.(lua.LString))
	ls.Desc = string(desc.(lua.LString))
	return nil
}
