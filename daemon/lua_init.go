package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

const luaScriptPath = "scripts"
const luaGlobalScriptName = "_script_name"

var luaValidScriptNames = []string{"main.lua", "init.lua"}
var luaFuncs map[string]lua.LGFunction = nil

func luaGetScriptPath() string {
	return luaScriptPath
}

func luaGetAvailableScripts() []string {
	dirList, _ := ioutil.ReadDir(luaGetScriptPath())
	out := make([]string, 0)
	for _, fileInfo := range dirList {
		name := fileInfo.Name()
		// single file script
		if !fileInfo.IsDir() {
			ext := filepath.Ext(name)
			if ext == ".lua" {
				out = append(out, strings.TrimSuffix(name, filepath.Ext(name)))
			}
			continue
		}
		// script inside directory
		for _, validScriptName := range luaValidScriptNames {
			pathTo := filepath.Join(luaGetScriptPath(), name, validScriptName)
			if _, err := os.Stat(pathTo); err == nil {
				out = append(out, name)
				break
			}
		}
	}
	return out
}

func luaGetEnabledScripts() []string {
	return []string{}
}

func luaGetPathToScript(name string) (string, error) {
	pathToFile := filepath.Join(luaGetScriptPath(), name+".lua")
	if _, err := os.Stat(pathToFile); err == nil {
		return pathToFile, nil
	}
	for _, validScriptName := range luaValidScriptNames {
		pathToFile = filepath.Join(luaGetScriptPath(), name, validScriptName)
		if _, err := os.Stat(pathToFile); err == nil {
			return pathToFile, nil
		}
	}
	return "", ErrLuaScriptNotFound
}

func luaRegisterFunction(name string, function lua.LGFunction) {
	if luaFuncs == nil {
		luaFuncs = make(map[string]lua.LGFunction)
	}
	luaFuncs[name] = function
}
