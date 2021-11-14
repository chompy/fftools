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

// DataTypeEncounter - Data type, encounter data
const DataTypeEncounter byte = 2

// Encounter - Data about an encounter
type Encounter struct {
	ByteEncodable
	ID           uint32    `json:"id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Zone         string    `json:"zone"`
	Damage       int32     `json:"damage"`
	Active       bool      `json:"active"`
	SuccessLevel uint8     `json:"success_level"`
}

// ToBytes - Convert to bytes
func (e *Encounter) ToBytes() []byte {
	data := make([]byte, 1)
	data[0] = DataTypeEncounter
	writeInt32(&data, int32(e.ID))
	writeTime(&data, e.StartTime)
	writeTime(&data, e.EndTime)
	writeString(&data, e.Zone)
	writeInt32(&data, e.Damage)
	writeBool(&data, e.Active)
	writeByte(&data, e.SuccessLevel)
	return data
}

// FromBytes - Convert act bytes to encounter
func (e *Encounter) FromBytes(data []byte) error {
	if data[0] != DataTypeEncounter {
		return ErrInvalidDataType
	}
	pos := 1
	e.ID = readUint32(data, &pos)
	e.StartTime = readTime(data, &pos)
	e.EndTime = readTime(data, &pos)
	e.Zone = readString(data, &pos)
	e.Damage = readInt32(data, &pos)
	e.Active = (readByte(data, &pos) != 0)
	e.SuccessLevel = readByte(data, &pos)
	return nil
}

func (e Encounter) ToLua() *lua.LTable {
	t := &lua.LTable{}
	t.RawSetString("id", lua.LNumber(e.ID))
	t.RawSetString("start_time", lua.LNumber(e.StartTime.Unix()))
	t.RawSetString("end_time", lua.LNumber(e.EndTime.Unix()))
	t.RawSetString("zone", lua.LString(e.Zone))
	t.RawSetString("damage", lua.LNumber(e.Damage))
	t.RawSetString("active", lua.LBool(e.Active))
	t.RawSetString("success_level", lua.LNumber(e.SuccessLevel))
	t.RawSetString("__goobject", &lua.LUserData{Value: e})
	return t
}
