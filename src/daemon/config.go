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
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configLuaDefault = "default.yaml"

var configLoadedApp *configApp
var proxyUid = ""
var proxySecret = ""

type configApp struct {
	PortData       uint16 `yaml:"data_port"`
	PortWeb        uint16 `yaml:"web_port"`
	LogMaxSize     int64  `yaml:"log_max_size"`
	EnableProxy    bool   `yaml:"enable_proxy"`
	ProxyAddress   string `yaml:"proxy_address"`
	ProxyURL       string `yaml:"proxy_url"`
	EnableKeyPress bool   `yaml:"enable_key_press"`
	EnableTTS      bool   `yaml:"enable_tts"`
}

func configGetPath() string {
	return filepath.Join(getBasePath(), configPath)
}

func configGetScriptsEnabledPath() string {
	return dataGetPath("_enabled")
}

func configGetProxyCredsPath() string {
	return dataGetPath("_proxy")
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

func configGetProxyCreds() (string, string, error) {
	if proxyUid != "" && proxySecret != "" {
		return proxyUid, proxySecret, nil
	}
	pathTo := configGetProxyCredsPath()
	raw, err := ioutil.ReadFile(pathTo)
	if err != nil {
		if os.IsNotExist(err) {
			uid, secret := webProxyGenerateCreds()
			return uid, secret, configSetProxyCred(uid, secret)
		}
		return "", "", err
	}
	creds := []string{}
	if err := json.Unmarshal(raw, &creds); err != nil {
		return "", "", err
	}
	proxyUid = creds[0]
	proxySecret = creds[1]
	return creds[0], creds[1], nil
}

func configSetProxyCred(uid string, secret string) error {
	proxyUid = uid
	proxySecret = secret
	pathTo := configGetProxyCredsPath()
	raw, err := json.Marshal([]string{uid, secret})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pathTo, raw, 0600)
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
		PortData:       31593,
		PortWeb:        31594,
		LogMaxSize:     262144, // 256KB
		EnableProxy:    true,
		EnableKeyPress: true,
		EnableTTS:      true,
		ProxyAddress:   "fftools.net:31595",
		//ProxyAddress: "localhost:31595",
		ProxyURL: "https://fftools.net/",
		//ProxyURL: "http://localhost:31596/",
	}
}

func configAppLoad() *configApp {
	if configLoadedApp != nil {
		return configLoadedApp
	}
	config := configAppDefault()
	rawConfig, err := ioutil.ReadFile(configGetScriptConfigPath("_app"))
	if err != nil {
		logWarn(err.Error())
		return config
	}
	if err := yaml.Unmarshal(rawConfig, config); err != nil {
		logWarn(err.Error())
	}
	configLoadedApp = config
	return config
}
