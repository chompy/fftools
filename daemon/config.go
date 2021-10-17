package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const configPath = "../config"
const configScriptsEnabledFile = "enabled.json"

func configGetPath() string {
	return configPath
}

func configGetScriptsEnabledPath() string {
	return filepath.Join(configGetPath(), configScriptsEnabledFile)
}

func configGetScriptConfigPath(name string) string {
	return filepath.Join(configGetPath(), name+".json")
}

func configLoadScriptsEnabled() ([]string, error) {
	pathTo := configGetScriptsEnabledPath()
	raw, err := ioutil.ReadFile(pathTo)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func configLoadScriptConfig(name string) (map[string]interface{}, error) {
	pathTo := configGetScriptConfigPath(name)
	raw, err := ioutil.ReadFile(pathTo)
	if err != nil {
		return nil, err
	}
	out := make(map[string]interface{})
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}
