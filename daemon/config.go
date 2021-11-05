package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configScriptsEnabledFile = "_enabled.json"
const configLuaDefault = "default.yaml"

var configLoadedApp *configApp

type configApp struct {
	PortData   uint16 `yaml:"port_data"`
	PortWeb    uint16 `yaml:"port_web"`
	LogMaxSize int64  `yaml:"log_max_size"`
}

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
	pathToFile := filepath.Join(getScriptPath(), name, configLuaDefault)
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

func configAppDefault() *configApp {
	return &configApp{
		PortData:   31593,
		PortWeb:    31594,
		LogMaxSize: 262144, // 256KB
	}
}

func configAppLoad() *configApp {
	if configLoadedApp != nil {
		return configLoadedApp
	}
	config := configAppDefault()
	rawConfig, err := ioutil.ReadFile(configGetScriptConfigPath("_app"))
	if err != nil {
		return config
	}
	if err := yaml.Unmarshal(rawConfig, config); err != nil {
		logWarn(err.Error())
	}
	configLoadedApp = config
	return config
}
