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
	"log"
	"runtime"
	"strconv"
	"strings"
)

// hexToInt converts hex string to an integer.
func hexToInt(hexString string) (int, error) {
	if hexString == "" {
		return 0, nil
	}
	output, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Println(err.Error(), hexString, fn, line)
	}
	return int(output), err
}

const roleTank = "tank"
const roleDps = "dps"
const roleHealer = "healer"

func jobGetRole(job string) string {
	switch strings.ToLower(job) {
	case "pld", "war", "drk", "gnb":
		{
			return roleTank
		}
	case "whm", "sch", "ast", "sge":
		{
			return roleHealer
		}
	}
	return roleDps
}
