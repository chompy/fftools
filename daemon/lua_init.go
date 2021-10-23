package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

const luaScriptPath = "../scripts"
const luaGlobalScriptName = "_script_name"

var luaValidScriptNames = []string{"main.lua", "init.lua"}
var luaFuncs map[string]lua.LGFunction = nil
var loadedScripts []*luaScript = make([]*luaScript, 0)

func luaGetScriptPath() string {
	return luaScriptPath
}

func luaGetAvailableScriptNames() []string {
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

func luaGetEnabledScriptNames() []string {
	available := luaGetAvailableScriptNames()
	enabled, err := configLoadScriptsEnabled()
	if err != nil {
		logWarn(err.Error())
		return []string{}
	}
	out := make([]string, 0)
	for _, enabledScript := range enabled {
		for _, availableScript := range available {
			if enabledScript == availableScript {
				out = append(out, availableScript)
				break
			}
		}
	}
	return out
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

func luaLoadScripts() []*luaScript {
	availableScriptNames := luaGetAvailableScriptNames()
	// remove scripts that are no longer available
	hasRemovedScript := true
	for hasRemovedScript {
		hasRemovedScript = false
		for i, script := range loadedScripts {
			hasScript := false
			for _, availableScriptName := range availableScriptNames {
				if availableScriptName == script.ScriptName {
					hasScript = true
					break
				}
			}
			if hasScript {
				continue
			}
			loadedScripts = append(loadedScripts[:i], loadedScripts[i+1:]...)
			hasRemovedScript = true
			break
		}
	}
	for _, scriptName := range availableScriptNames {
		// check if loaded already
		scriptIsLoaded := false
		for _, loadedScript := range loadedScripts {
			if loadedScript.ScriptName == scriptName {
				scriptIsLoaded = true
				break
			}
		}
		if scriptIsLoaded {
			continue
		}
		// load script
		luaScript, _ := luaLoadScript(scriptName)
		loadedScripts = append(loadedScripts, luaScript)
	}
	return loadedScripts
}

func luaEnableScripts() {
	enabledScriptNames := luaGetEnabledScriptNames()
	for _, script := range luaLoadScripts() {
		// check enabled
		enabled := false
		for _, enabledScriptName := range enabledScriptNames {
			if enabledScriptName == script.ScriptName {
				enabled = true
				break
			}
		}
		if enabled && !script.Enabled {
			// previously disabled, now enabled
			if script.L == nil {
				if err := script.load(); err != nil {
					logLuaWarn(script.L, err.Error())
					actError(err, script.ScriptName)
					continue
				}
			}
			if err := script.init(); err != nil {
				logLuaWarn(script.L, err.Error())
				actError(err, script.ScriptName)
				continue
			}
			script.Enabled = true
		} else if !enabled && script.Enabled {
			// previously enabled, now disabled
			script.close()
			script.Enabled = false
		}
	}
}
