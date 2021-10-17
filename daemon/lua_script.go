package main

import lua "github.com/yuin/gopher-lua"

type luaScript struct {
	Name string
	Path string
	L    *lua.LState
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
		Name: name,
		Path: pathTo,
		L:    L,
	}
	return ls, nil
}

func (ls *luaScript) close() {
	ls.L.Close()
}

func (ls *luaScript) init() error {
	logInfo("Init '%s' Lua script.", ls.Name)
	return ls.L.DoString("init()")
}
