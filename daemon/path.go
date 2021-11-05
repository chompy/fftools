package main

import (
	"os"
	"path/filepath"
)

const configPath = "config"
const dataPath = "data"
const scriptPath = "scripts"
const scriptWebPath = "web"
const logPath = "log"

func getBasePath() string {
	exePath, _ := os.Executable()
	return filepath.Dir(exePath)
}

func getScriptPath() string {
	return filepath.Join(getBasePath(), scriptPath)
}

func getScriptWebPath(name string) string {
	return filepath.Join(getScriptPath(), name, scriptWebPath)
}

func getLogPath() string {
	return filepath.Join(getBasePath(), logPath)
}
