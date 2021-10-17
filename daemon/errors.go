package main

import (
	"errors"
)

var (
	ErrInvalidDataType   = errors.New("invalid data type")
	ErrLuaScriptNotFound = errors.New("lua script not found")
)
