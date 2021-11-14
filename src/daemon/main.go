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
)

func main() {
	luaEnableScripts()
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		testLogLines := testerParse(os.Stdin)
		if len(testLogLines) > 0 {
			testerReplay(testLogLines)
		}
	}
	config := configAppLoad()
	if config.EnableProxy {
		go webProxyConnect()
	}
	go initWeb()
	actListenUDP()
}
