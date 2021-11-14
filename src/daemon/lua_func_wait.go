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
	"time"

	lua "github.com/yuin/gopher-lua"
)

func luaFuncWait(L *lua.LState) int {
	waitTime := L.ToInt(1)
	callback := L.ToFunction(2)
	if callback == nil {
		logLuaWarn(L, "Invalid callback for 'wait' function.")
		return 0
	}
	go func() {
		time.Sleep(time.Millisecond * time.Duration(waitTime))
		L.Push(callback)
		if err := L.PCall(0, 0, nil); err != nil {
			logLuaWarn(L, err.Error())
		}
	}()
	return 0
}

func init() {
	luaRegisterFunction("wait", luaFuncWait)
}
