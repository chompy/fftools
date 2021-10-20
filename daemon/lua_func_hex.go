package main

import (
	"fmt"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func luaFuncIntToHex(L *lua.LState) int {
	num := L.ToInt64(1)
	L.Push(lua.LString(strings.ToUpper(fmt.Sprintf("%x", num))))
	return 1
}

func luaFuncHexToInt(L *lua.LState) int {
	hex := L.ToString(1)
	num, _ := strconv.ParseInt(hex, 16, 64)
	L.Push(lua.LNumber(num))
	return 1
}

func init() {
	luaRegisterFunction("int_to_hex", luaFuncIntToHex)
	luaRegisterFunction("hex_to_int", luaFuncHexToInt)

}
