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
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	lua "github.com/yuin/gopher-lua"
)

var keyPressQueue = make([][]string, 0)

func processKeyPressQueue() {
	for {
		for i, keys := range keyPressQueue {
			for _, key := range keys {
				robotgo.KeyDown(key)
				robotgo.MilliSleep(25)
			}
			for _, key := range keys {
				robotgo.KeyUp(key)
			}
			robotgo.MilliSleep(25)
			keyPressQueue = append(keyPressQueue[:i], keyPressQueue[i+1:]...)
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func luaFuncKeyPress(L *lua.LState) int {
	keys := make([]string, 0)
	index := 1
	for {
		v := L.Get(index)
		if v.Type() != lua.LTString {
			break
		}
		keys = append(keys, string(v.(lua.LString)))
		index++
	}
	logLuaInfo(L, "Press key '%s'.", strings.Join(keys, "-"))
	keyPressQueue = append(keyPressQueue, keys)
	return 0
}

func init() {
	luaRegisterFunction("key_press", luaFuncKeyPress)
	go processKeyPressQueue()
}
