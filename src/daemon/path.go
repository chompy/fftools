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
	"os"
	"path/filepath"
)

const configPath = "config"
const dataPath = "data"
const scriptPath = "scripts"
const scriptWebPath = "web"

func getBasePath() string {
	exePath, _ := os.Executable()
	return filepath.Join(filepath.Dir(exePath), "..")
}

func getScriptPath() string {
	return filepath.Join(getBasePath(), scriptPath)
}

func getScriptWebPath(name string) string {
	return filepath.Join(getScriptPath(), name, scriptWebPath)
}

/*func getLogPath() string {
	return filepath.Join(getBasePath(), logPath)
}*/
