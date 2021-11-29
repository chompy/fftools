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
	"log"
	"os"
	"path/filepath"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"gopkg.in/natefinch/lumberjack.v2"
)

const version = "0.06"

func main() {
	// set log output
	log.SetOutput(&lumberjack.Logger{
		Filename:   filepath.Join(getBasePath(), logPath),
		MaxSize:    10,
		MaxBackups: 3,
	})
	// set auto update
	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/chompy/fftools",
			ArchiveName:   "fftools_update.zip",
		},
		ExecutableName: "fftools_daemon.exe",
		Version:        "v" + version,
	}
	go func() {
		s, err := u.Update()
		if err != nil {
			logWarn("[UPDATE] " + err.Error())
			return
		}
		if s == updater.Updated {
			actRequestUpdate()
		}
	}()

	// enable scripts
	luaEnableScripts()
	// run testers if test data provided
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		testLogLines := testerParse(os.Stdin)
		if len(testLogLines) > 0 {
			testerReplay(testLogLines)
		}
	}
	// enable proxy
	config := configAppLoad()
	if config.EnableProxy {
		go webProxyConnect()
	}
	// enable web view
	go initWeb()
	// listen for act packets
	if err := actListenUDP(); err != nil {
		logPanic(err)
	}
}
