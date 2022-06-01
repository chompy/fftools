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

// DataTypeLogLine - Data type, log line
const DataTypeLogLine byte = 5

// LogLine - Log line from Act
type LogLine struct {
	ByteEncodable
	EncounterID uint32
	Time        time.Time
	LogLine     string
}

// ToBytes - Convert to bytes
func (l *LogLine) ToBytes() []byte {
	data := make([]byte, 1)
	data[0] = DataTypeLogLine
	writeInt32(&data, int32(l.EncounterID))
	writeTime(&data, l.Time)
	writeString(&data, l.LogLine)
	return data
}

// FromBytes - Convert bytes to log line
func (l *LogLine) FromBytes(data []byte) error {
	if data[0] != DataTypeLogLine {
		return ErrInvalidDataType
	}
	pos := 1
	l.EncounterID = readUint32(data, &pos)
	l.Time = readTime(data, &pos)
	l.LogLine = readString(data, &pos)
	return nil
}

func (l LogLine) ToLua() *lua.LTable {
	t := &lua.LTable{}
	t.RawSetString("encounter_id", lua.LNumber(l.EncounterID))
	t.RawSetString("time", lua.LNumber(l.Time.Unix()))
	t.RawSetString("log_line", lua.LString(l.LogLine))
	t.RawSetString("__goobject", &lua.LUserData{Value: l})
	return t
}
