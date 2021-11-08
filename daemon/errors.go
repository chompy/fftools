/*
This file is part of FF Lua.

FF Lua is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

FF Lua is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with FF Lua.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
)

var (
	ErrInvalidDataType             = errors.New("invalid data type")
	ErrLuaScriptNotFound           = errors.New("lua script not found")
	ErrDefaultConfigNotFound       = errors.New("default config file not found")
	ErrLuaScriptNotLoaded          = errors.New("lua script not loaded")
	ErrLuaScriptInError            = errors.New("lua script in error state")
	ErrConfigNotFound              = errors.New("config file not found")
	ErrActNotConnected             = errors.New("no act connection available")
	ErrLogParseTooFewCharacters    = errors.New("tried to parse log line with too few characters")
	ErrLogParseAbilityTooFewFields = errors.New("not enough fields when parsing ability")
)
