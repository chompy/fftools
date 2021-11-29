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

	"github.com/micmonay/keybd_event"
	lua "github.com/yuin/gopher-lua"
)

var keyPressQueue = make([][]string, 0)

var keyMapping = map[string]int{
	"f1":        keybd_event.VK_F1,
	"f2":        keybd_event.VK_F2,
	"f3":        keybd_event.VK_F3,
	"f4":        keybd_event.VK_F4,
	"f5":        keybd_event.VK_F5,
	"f6":        keybd_event.VK_F6,
	"f7":        keybd_event.VK_F7,
	"f8":        keybd_event.VK_F8,
	"f9":        keybd_event.VK_F9,
	"f10":       keybd_event.VK_F10,
	"f11":       keybd_event.VK_F11,
	"f12":       keybd_event.VK_F12,
	"1":         keybd_event.VK_1,
	"2":         keybd_event.VK_2,
	"3":         keybd_event.VK_3,
	"4":         keybd_event.VK_4,
	"5":         keybd_event.VK_5,
	"6":         keybd_event.VK_6,
	"7":         keybd_event.VK_7,
	"8":         keybd_event.VK_8,
	"9":         keybd_event.VK_9,
	"0":         keybd_event.VK_0,
	"minus":     keybd_event.VK_MINUS,
	"=":         keybd_event.VK_EQUAL,
	"backspace": keybd_event.VK_BACKSPACE,
	"tab":       keybd_event.VK_TAB,
	"[":         keybd_event.VK_LEFTBRACE,
	"]":         keybd_event.VK_RIGHTBRACE,
	"enter":     keybd_event.VK_ENTER,
	";":         keybd_event.VK_SEMICOLON,
	"'":         keybd_event.VK_APOSTROPHE,
	"\\":        keybd_event.VK_BACKSLASH,
	",":         keybd_event.VK_COMMA,
	".":         keybd_event.VK_DOT,
	"/":         keybd_event.VK_SLASH,
	"caps":      keybd_event.VK_CAPSLOCK,
	"q":         keybd_event.VK_Q,
	"w":         keybd_event.VK_W,
	"e":         keybd_event.VK_E,
	"r":         keybd_event.VK_R,
	"t":         keybd_event.VK_T,
	"y":         keybd_event.VK_Y,
	"u":         keybd_event.VK_U,
	"i":         keybd_event.VK_I,
	"o":         keybd_event.VK_O,
	"p":         keybd_event.VK_P,
	"a":         keybd_event.VK_A,
	"s":         keybd_event.VK_S,
	"d":         keybd_event.VK_D,
	"f":         keybd_event.VK_F,
	"g":         keybd_event.VK_G,
	"h":         keybd_event.VK_H,
	"j":         keybd_event.VK_J,
	"k":         keybd_event.VK_K,
	"l":         keybd_event.VK_L,
	"z":         keybd_event.VK_Z,
	"x":         keybd_event.VK_X,
	"c":         keybd_event.VK_C,
	"v":         keybd_event.VK_V,
	"b":         keybd_event.VK_B,
	"n":         keybd_event.VK_N,
	"m":         keybd_event.VK_M,
}

func processKeyPressQueue() {
	config := configAppLoad()
	if !config.EnableKeyPress {
		logInfo("Key presses are disabled.")
		return
	}
	logInfo("Key presses are enabled.")
	for {
		for i, keys := range keyPressQueue {
			kb, err := keybd_event.NewKeyBonding()
			if err != nil {
				logWarn(err.Error())
				continue
			}
			for _, key := range keys {
				switch key {
				case "alt":
					{
						kb.HasALT(true)
						break
					}
				case "ctrl":
					{
						kb.HasCTRL(true)
						break
					}
				case "shift":
					{
						kb.HasSHIFT(true)
						break
					}
				default:
					{
						if keyMapping[key] != 0 {
							kb.AddKey(keyMapping[key])
						}
						break
					}
				}
			}
			kb.Press()
			time.Sleep(time.Millisecond * 15)
			kb.Release()
			keyPressQueue = append(keyPressQueue[:i], keyPressQueue[i+1:]...)
			break
		}
		time.Sleep(time.Millisecond * 250)
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
