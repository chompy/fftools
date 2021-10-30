package main

import (
	"os"
	"path/filepath"
)

const configPath = "config"
const dataPath = "data"
const scriptPath = "scripts"

func getBasePath() string {
	exePath, _ := os.Executable()
	return filepath.Dir(exePath)
}
