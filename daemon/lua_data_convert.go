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

func valueLuaTableToGo(value *lua.LTable) map[interface{}]interface{} {
	out := make(map[interface{}]interface{})
	value.ForEach(func(k lua.LValue, v lua.LValue) {
		out[valueLuaToGo(k)] = valueLuaToGo(v)
	})
	return out
}
