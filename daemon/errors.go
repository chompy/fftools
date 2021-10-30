package main

import (
	"errors"
)

var (
	ErrInvalidDataType             = errors.New("invalid data type")
	ErrLuaScriptNotFound           = errors.New("lua script not found")
	ErrDefaultConfigNotFound       = errors.New("default config file not found")
	ErrConfigNotFound              = errors.New("config file not found")
	ErrActNotConnected             = errors.New("no act connection available")
	ErrLogParseTooFewCharacters    = errors.New("tried to parse log line with too few characters")
	ErrLogParseAbilityTooFewFields = errors.New("not enough fields when parsing ability")
)
