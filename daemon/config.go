package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configScriptsEnabledFile = "enabled.json"
const luaDefaultConfigFile = "default.yaml"

func configGetPath() string {
	return filepath.Join(getBasePath(), configPath)
}

func configGetScriptsEnabledPath() string {
	return filepath.Join(configGetPath(), configScriptsEnabledFile)
}

func configGetScriptConfigPath(name string) string {
	return filepath.Join(configGetPath(), name+".yaml")
}

func configLoadScriptsEnabled() ([]string, error) {
	pathTo := configGetScriptsEnabledPath()
	raw, err := ioutil.ReadFile(pathTo)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	out := make([]string, 0)
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func configSetScriptEnabled(name string, enable bool) error {
	enabledScripts, err := configLoadScriptsEnabled()
	if err != nil {
		return err
	}
	enabledIndex := -1
	for index, enabledScript := range enabledScripts {
		if enabledScript == name {
			enabledIndex = index
			break
		}
	}
	if !enable && enabledIndex > -1 {
		enabledScripts = append(enabledScripts[:enabledIndex], enabledScripts[enabledIndex+1:]...)
	} else if enable && enabledIndex == -1 {
		enabledScripts = append(enabledScripts, name)
	}
	enabledScriptsJson, err := json.Marshal(enabledScripts)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configGetScriptsEnabledPath(), enabledScriptsJson, 0755)
}

func configGetPathToScriptDefaultConfig(name string) (string, error) {
	pathToFile := filepath.Join(getScriptPath(), name, luaDefaultConfigFile)
	if _, err := os.Stat(pathToFile); err != nil {
		if os.IsNotExist(err) {
			return "", ErrDefaultConfigNotFound
		}
		return "", err
	}
	return pathToFile, nil
}

func configLoadScriptConfig(name string) (map[string]interface{}, error) {
	pathTo := configGetScriptConfigPath(name)
	raw, err := ioutil.ReadFile(pathTo)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// copy default config file
		defaultPath, err := configGetPathToScriptDefaultConfig(name)
		if err != nil {
			return nil, err
		}
		raw, err = ioutil.ReadFile(defaultPath)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(pathTo, raw, 0755); err != nil {
			return nil, err
		}
		logInfo("Created default config for %s.", name)
	}
	out := make(map[string]interface{})
	if err := yaml.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}
