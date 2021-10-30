/*
This file is part of FFLiveParse.

FFLiveParse is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

FFLiveParse is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with FFLiveParse.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

// DataTypeCombatant - Data type, combatant data
const DataTypeCombatant byte = 3

// Combatant - Data about a combatant
type Combatant struct {
	ByteEncodable
	ID           int32     `json:"id"`
	Name         string    `json:"name"`
	EncounterID  uint32    `json:"act_encounter_id"`
	Time         time.Time `json:"time"`
	Job          string    `json:"job"`
	Damage       int32     `json:"damage"`
	DamageTaken  int32     `json:"damage_taken"`
	DamageHealed int32     `json:"damage_healed"`
	Deaths       int32     `json:"deaths"`
	Hits         int32     `json:"hits"`
	Heals        int32     `json:"heals"`
	Kills        int32     `json:"kills"`
	CritHits     int32     `json:"critical_hits"`
	CritHeals    int32     `json:"critical_heals"`
}

// ToBytes - Convert to bytes
func (c *Combatant) ToBytes() []byte {
	data := make([]byte, 1)
	data[0] = DataTypeCombatant
	writeInt32(&data, int32(c.EncounterID))
	writeInt32(&data, c.ID)
	writeString(&data, c.Name)
	writeString(&data, c.Job)
	writeInt32(&data, c.Damage)
	writeInt32(&data, c.DamageTaken)
	writeInt32(&data, c.DamageHealed)
	writeInt32(&data, c.Deaths)
	writeInt32(&data, c.Hits)
	writeInt32(&data, c.Heals)
	writeInt32(&data, c.Kills)
	writeInt32(&data, c.CritHits)
	writeInt32(&data, c.CritHeals)
	writeTime(&data, c.Time)
	return data
}

// FromBytes - Convert act bytes to combatant
func (c *Combatant) FromBytes(data []byte) error {
	if data[0] != DataTypeCombatant {
		return ErrInvalidDataType
	}
	pos := 1
	c.EncounterID = readUint32(data, &pos)
	c.ID = readInt32(data, &pos)
	c.Name = readString(data, &pos)
	c.Job = readString(data, &pos)
	c.Damage = readInt32(data, &pos)
	c.DamageTaken = readInt32(data, &pos)
	c.DamageHealed = readInt32(data, &pos)
	c.Deaths = readInt32(data, &pos)
	c.Hits = readInt32(data, &pos)
	c.Heals = readInt32(data, &pos)
	c.Kills = readInt32(data, &pos)
	c.CritHits = readInt32(data, &pos)
	c.CritHeals = readInt32(data, &pos)
	c.Time = time.Now()
	return nil
}

func (c Combatant) ToLua() *lua.LTable {
	t := &lua.LTable{}
	t.RawSetString("encounter_id", lua.LNumber(c.EncounterID))
	t.RawSetString("time", lua.LNumber(c.Time.Unix()))
	t.RawSetString("id", lua.LNumber(c.ID))
	t.RawSetString("name", lua.LString(c.Name))
	t.RawSetString("job", lua.LString(c.Job))
	t.RawSetString("damage", lua.LNumber(c.Damage))
	t.RawSetString("damage_taken", lua.LNumber(c.DamageTaken))
	t.RawSetString("damage_healed", lua.LNumber(c.DamageHealed))
	t.RawSetString("deaths", lua.LNumber(c.Deaths))
	t.RawSetString("hits", lua.LNumber(c.Hits))
	t.RawSetString("heals", lua.LNumber(c.Heals))
	t.RawSetString("kills", lua.LNumber(c.Kills))
	t.RawSetString("critical_hits", lua.LNumber(c.CritHits))
	t.RawSetString("critical_heals", lua.LNumber(c.CritHeals))
	t.RawSetString("__goobject", &lua.LUserData{Value: c})
	return t
}
