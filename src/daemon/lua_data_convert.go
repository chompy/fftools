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

import lua "github.com/yuin/gopher-lua"

func valueGoToLua(value interface{}) lua.LValue {
	switch value := value.(type) {
	case string:
		{
			return lua.LString(value)
		}
	case int:
		{
			return lua.LNumber(value)
		}
	case float64:
		{
			return lua.LNumber(value)
		}
	case bool:
		{
			return lua.LBool(value)
		}
	case map[string]interface{}:
		{
			return valueGoToLuaTable(value)
		}
	case map[interface{}]interface{}:
		{
			return valueGoToLuaTable(value)
		}
	case []interface{}:
		{
			return valueGoToLuaTable(value)
		}
	default:
		{
			return lua.LNil
		}
	}
}

func valueGoToLuaTable(value interface{}) *lua.LTable {
	switch value := value.(type) {
	case map[string]interface{}:
		{
			out := &lua.LTable{}
			for k, v := range value {
				out.RawSetString(k, valueGoToLua(v))
			}
			return out
		}
	case map[interface{}]interface{}:
		{
			out := &lua.LTable{}
			for k, v := range value {
				out.RawSet(valueGoToLua(k), valueGoToLua(v))
			}
			return out
		}
	case []interface{}:
		{
			out := &lua.LTable{}
			for k, v := range value {
				out.RawSetInt(k+1, valueGoToLua(v))
			}
			return out
		}
	default:
		{
			return nil
		}
	}
}

func valueLuaToGo(value lua.LValue) interface{} {
	switch value.Type() {
	case lua.LTString:
		{
			return string(value.(lua.LString))
		}
	case lua.LTNumber:
		{
			return float64(value.(lua.LNumber))
		}
	case lua.LTBool:
		{
			return bool(value.(lua.LBool))
		}
	case lua.LTTable:
		{
			return valueLuaTableToGo(value.(*lua.LTable))
		}
	default:
		{
			return nil
		}
	}
}

func valueLuaTableToGo(value *lua.LTable) map[string]interface{} {
	out := make(map[string]interface{})
	value.ForEach(func(k lua.LValue, v lua.LValue) {
		out[k.String()] = valueLuaToGo(v)
	})
	return out
}
