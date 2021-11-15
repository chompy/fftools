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
	"regexp"
	"strconv"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// LogTypeGameLog defines a game log event.
const LogTypeGameLog = 0x00

// LogTypeZoneChange defines a zone change event.
const LogTypeZoneChange = 0x01

// LogTypeChangePrimaryPlayer defines a player change event.
const LogTypeChangePrimaryPlayer = 0x02

// LogTypeAddPlayer defines an add combatant event.
const LogTypeAddCombatant = 0x03

// LogTypeRemoveCombatant defines a remove combatant event.
const LogTypeRemoveCombatant = 0x04

// LogTypeAddBuff defines an add buff event.
const LogTypeAddBuff = 0x05

// LogTypeRemoveBuff defines a remove buff event.
const LogTypeRemoveBuff = 0x06

// LogTypeFlyingText defines a flying text event.
const LogTypeFlyingText = 0x07

// LogTypeOutgoingAbility defines an outgoing ability event.
const LogTypeOutgoingAbility = 0x08

// LogTypeIncomingAbility defines an incoming ability event.
const LogTypeIncomingAbility = 0x09

// LogTypePartyList defines a party list event (party makeup changes).
const LogTypePartyList = 0x0B

// LogTypePlayerStats defines a player stat event.
const LogTypePlayerStats = 0x0C

// LogTypeCombatantHP defines a combatant hp event.
const LogTypeCombatantHP = 0x0D

// LogTypeNetworkStartsCasting defines a network start casting event.
const LogTypeNetworkStartsCasting = 0x14

// LogTypeNetworkAbility defines a network ability event.
const LogTypeNetworkAbility = 0x15

// LogTypeNetworkAOEAbility defines a network aoe ability event.
const LogTypeNetworkAOEAbility = 0x16

// LogTypeNetworkCancelAbility defines a network cancel ability event.
const LogTypeNetworkCancelAbility = 0x17

// LogTypeNetworkDot defines a network dot event.
const LogTypeNetworkDot = 0x18

// LogTypeNetworkDeath defines a network death event.
const LogTypeNetworkDeath = 0x19

// LogTypeNetworkBuff defines a network buff event.
const LogTypeNetworkBuff = 0x1A

// LogTypeNetworkTargetIcon defines a network target icon event.
const LogTypeNetworkTargetIcon = 0x1B

// LogTypeNetworkRaidMarker defines a network raid marker event.
const LogTypeNetworkRaidMarker = 0x1C

// LogTypeNetworkTargetMarker defines a network target marker event.
const LogTypeNetworkTargetMarker = 0x1D

// LogTypeNetworkBuffRemove defines a network buff remove event.
const LogTypeNetworkBuffRemove = 0x1E

// LogTypeNetworkGauge defines a network guage event.
const LogTypeNetworkGauge = 0x1F

// LogTypeNetwork6D defines a network 6D event.
const LogTypeNetwork6D = 0x21

// LogTypeNetworkNameToggle defines a network name toggle event.
const LogTypeNetworkNameToggle = 0x22

// LogTypeNetworkTether defines a network tether event (tether between two combatants).
const LogTypeNetworkTether = 0x23

// LogTypeLimitBreak defines a limit break event.
const LogTypeLimitBreak = 0x24

// LogTypeNetworkActionSync defines a network action sync event.
const LogTypeNetworkActionSync = 0x25

// LogTypeNetworkStatusEffects defines a network status effects event.
const LogTypeNetworkStatusEffects = 0x26

// LogTypeNetworkUpdateHP defines a network up hp event.
const LogTypeNetworkUpdateHP = 0x27

// LogTypeMap defines a map event (map change).
const LogTypeMap = 0x28

// LogFieldType - Log field identifier, message type
const LogFieldType = 0

// LogFieldAttackerID - Log field identifier, attacker id
const LogFieldAttackerID = 1

// LogFieldAttackerName - Log field identifier, attacker name
const LogFieldAttackerName = 2

// LogFieldAbilityID - Log field identifier, ability id
const LogFieldAbilityID = 3

// LogFieldAbilityName - Log field identifier, ability name
const LogFieldAbilityName = 4

// LogFieldTargetID - Log field identifier, target id
const LogFieldTargetID = 5

// LogFieldTargetName - Log field identifier, target name
const LogFieldTargetName = 6

// LogFieldFlags - Log field identifier, flags
const LogFieldFlags = 7

// LogFieldDamage - Log field identifier, damage
const LogFieldDamage = 8

// LogFieldTargetCurrentHP - Log field identifier, target current hp
const LogFieldTargetCurrentHP = 23

// LogFieldTargetMaxHP - Log field identifier, target max hp
const LogFieldTargetMaxHP = 24

// LogFieldTargetX - Log field identifier, target x pos
const LogFieldTargetX = 30

// LogFieldTargetY - Log field identifier, target y pos
const LogFieldTargetY = 31

// LogFieldTargetZ - Log field identifier, target z pos
const LogFieldTargetZ = 32

// LogFieldAttackerCurrentHP - Log field identifier, attacker current hp
const LogFieldAttackerCurrentHP = 33

// LogFieldAttackerMaxHP - Log field identifier, attacker max hp
const LogFieldAttackerMaxHP = 34

// LogFieldAttackerX - Log field identifier, attacker x pos
const LogFieldAttackerX = 39

// LogFieldAttackerY - Log field identifier, attacker y pos
const LogFieldAttackerY = 40

// LogFieldAttackerZ - Log field identifier, attacker z pos
const LogFieldAttackerZ = 41

// logShiftValues
var logShiftValues = [...]int{0x3E, 0x113, 0x213, 0x313}

// logRegexes is a map of regular expressions used to parse log events.
var logRegexes = map[int]*parseInstructions{
	LogTypeZoneChange: &parseInstructions{
		regexp.MustCompile(` 01:Changed Zone to (.*)\.`),
		[]string{"zone|str"},
	},
	LogTypeChangePrimaryPlayer: &parseInstructions{
		regexp.MustCompile(` 02:Changed primary player to (.*)\.`),
		[]string{"name|str"},
	},
	LogTypeNetworkDeath: &parseInstructions{
		regexp.MustCompile(` 19:([a-zA-Z0-9'\- ]*) was defeated by ([a-zA-Z0-9'\- ]*)`),
		[]string{"target_name|str", "source_name|str"},
	},
	LogTypeRemoveCombatant: &parseInstructions{
		regexp.MustCompile(` 04:([A-F0-9]*):Removing combatant ([a-zA-Z0-9'\- ]*)\.  Max HP####([0-9]*)\.`),
		[]string{"target_id|int", "target_name|str", "target_max_hp|int"},
	},
	LogTypeNetworkTargetIcon: &parseInstructions{
		regexp.MustCompile(` 1B:([A-F0-9]*):([a-zA-Z0-9'\- ]*):....:....:([A-F0-9]*):`),
		[]string{"target_id|int", "target_name|str", "icon_id|int"},
	},
	LogTypeNetworkBuff: &parseInstructions{
		regexp.MustCompile(` 1A:([A-F0-9]*):([a-zA-Z0-9'\- ]*) gains the effect of (.*) from (.*) for (.*) Seconds`),
		[]string{"target_id|int", "target_name|str", "effect_name|str", "effect_duration|float"},
	},
}

// parseInstructions defines how to parse a given log type.
type parseInstructions struct {
	Regexp *regexp.Regexp
	Fields []string
}

// ParsedLogEvent defines a log event that has been parsed.
type ParsedLogEvent struct {
	Type   int
	Raw    string
	Time   time.Time
	Values map[string]interface{}
}

// ParseLogEvent parses a log event and returns results as a ParsedLogEvent.
func ParseLogEvent(logLine LogLine) (ParsedLogEvent, error) {
	logLineString := logLine.LogLine
	if len(logLineString) <= 17 {
		return ParsedLogEvent{}, ErrLogParseTooFewCharacters
	}
	// get field type
	logLineType, err := hexToInt(logLineString[15:17])
	if err != nil {
		return ParsedLogEvent{}, err
	}
	// semi colon with space afterwards is ability name instead of delimiter
	// probably........... examples... Kaeshi: Higanbana, Hissatsu: Guren
	logLineString = strings.Replace(logLineString, ": ", "####", -1)
	// split fields
	fields := strings.Split(logLineString[15:], ":")
	// create data object
	out := ParsedLogEvent{
		Type:   int(logLineType),
		Raw:    strings.Replace(logLineString, "####", ": ", -1),
		Time:   logLine.Time,
		Values: make(map[string]interface{}),
	}
	// parse remaining
	switch logLineType {
	case LogTypeNetworkAbility, LogTypeNetworkAOEAbility:
		{
			// ensure there are enough fields
			if len(fields) < 40 {
				return out, ErrLogParseAbilityTooFewFields
			}
			// Shift damage and flags forward for mysterious spurious :3E:0:.
			// Plenary Indulgence also appears to prepend confession stacks.
			// UNKNOWN: Can these two happen at the same time?
			flagsInt, err := hexToInt(fields[LogFieldFlags])
			if err != nil {
				return out, err
			}
			for _, shiftValue := range logShiftValues {
				if flagsInt == shiftValue {
					fields[LogFieldFlags] = fields[LogFieldFlags+2]
					fields[LogFieldFlags+1] = fields[LogFieldFlags+3]
					break
				}
			}
			// fetch damage value
			damageFieldLength := len(fields[LogFieldDamage])
			damage := 0
			if damageFieldLength >= 4 {
				// Get the left four bytes as damage.
				damage, err = hexToInt(fields[LogFieldDamage][0:4])
				if err != nil {
					return out, err
				}
			}
			// Check for third byte == 0x40.
			if damageFieldLength >= 4 && fields[LogFieldDamage][damageFieldLength-4] == '4' {
				// Wrap in the 4th byte as extra damage.  See notes above.
				rightDamage, err := hexToInt(fields[LogFieldDamage][damageFieldLength-2 : damageFieldLength])
				if err != nil {
					return out, err
				}
				damage = damage - rightDamage + (rightDamage << 16)
			}
			out.Values["damage"] = int(damage)
			out.Values["source_id"], _ = hexToInt(fields[LogFieldAttackerID])
			out.Values["source_name"] = fields[LogFieldAttackerName]
			out.Values["ability_id"], _ = hexToInt(fields[LogFieldAbilityID])
			out.Values["ability_name"] = fields[LogFieldAbilityName]
			out.Values["ability_name"] = strings.Replace(out.Values["ability_name"].(string), "####", ": ", -1)
			out.Values["target_id"], _ = hexToInt(fields[LogFieldTargetID])
			out.Values["target_name"] = fields[LogFieldTargetName]
			out.Values["target_current_hp"], _ = strconv.ParseInt(fields[LogFieldTargetCurrentHP], 10, 64)
			out.Values["target_max_hp"], _ = strconv.ParseInt(fields[LogFieldTargetMaxHP], 10, 64)
			out.Values["target_x"], _ = strconv.ParseFloat(fields[LogFieldTargetX], 64)
			out.Values["target_y"], _ = strconv.ParseFloat(fields[LogFieldTargetY], 64)
			out.Values["target_z"], _ = strconv.ParseFloat(fields[LogFieldTargetZ], 64)
			out.Values["source_current_hp"], _ = strconv.ParseInt(fields[LogFieldAttackerCurrentHP], 10, 64)
			out.Values["source_max_hp"], _ = strconv.ParseInt(fields[LogFieldAttackerMaxHP], 10, 64)
			out.Values["source_x"], _ = strconv.ParseFloat(fields[LogFieldAttackerX], 64)
			out.Values["source_y"], _ = strconv.ParseFloat(fields[LogFieldAttackerY], 64)
			out.Values["source_z"], _ = strconv.ParseFloat(fields[LogFieldAttackerZ], 64)
			// flags
			out.Values["flag_dodge"] = false
			out.Values["flag_instant_death"] = false
			out.Values["flag_damage"] = false
			out.Values["flag_critical_hit"] = false
			out.Values["flag_direct_hit"] = false
			out.Values["flag_heal"] = false
			out.Values["flag_block"] = false
			out.Values["flag_parry"] = false
			if len(fields) >= LogFieldFlags+1 {
				rawFlags := fields[LogFieldFlags]
				if len(rawFlags) > 0 {
					switch rawFlags[len(rawFlags)-1:] {
					case "1":
						{
							out.Values["flag_dodge"] = true
							break
						}
					case "3":
						{
							if len(rawFlags) >= 4 {
								switch rawFlags[len(rawFlags)-3 : len(rawFlags)-2] {
								case "3":
									{
										out.Values["flag_instant_death"] = true
										break
									}
								default:
									{
										out.Values["flag_damage"] = true
										switch rawFlags[len(rawFlags)-4 : len(rawFlags)-3] {
										case "1":
											{
												out.Values["flag_critical_hit"] = true
												break
											}
										case "2":
											{
												out.Values["flag_direct_hit"] = true
												break
											}
										case "3":
											{
												out.Values["flag_critical_hit"] = true
												out.Values["flag_direct_hit"] = true
												break
											}
										}
										break
									}
								}
							}
							break
						}
					case "4":
						{
							out.Values["flag_heal"] = true
							if len(rawFlags) >= 6 && rawFlags[len(rawFlags)-6:len(rawFlags)-5] == "1" {
								out.Values["flag_critical_hit"] = true
							}
							break
						}
					case "5":
						{
							out.Values["flag_block"] = true
							break
						}
					case "6":
						{
							out.Values["flag_parry"] = true
							break
						}
					}
				}
			}
			break
		}
	case LogTypePartyList:
		{
			count, _ := strconv.ParseInt(fields[1], 10, 64)
			out.Values["count"] = count
			out.Values["ids"] = make([]int, count)
			for i := 0; i < int(count); i++ {
				out.Values["ids"].([]int)[i], _ = hexToInt(fields[i+2])
			}
			break
		}
	default:
		{
			instructions := logRegexes[out.Type]
			if instructions != nil {
				match := instructions.Regexp.FindAllStringSubmatch(out.Raw, -1)
				for index, field := range instructions.Fields {
					fieldSplit := strings.Split(field, "|")
					switch fieldSplit[1] {
					case "int":
						{
							out.Values[fieldSplit[0]] = 0
							if len(match) > 0 {
								value, err := hexToInt(match[0][index+1])
								if err == nil {
									out.Values[fieldSplit[0]] = value
								}
							}
							break
						}
					case "str":
						{
							out.Values[fieldSplit[0]] = ""
							if len(match) > 0 {
								out.Values[fieldSplit[0]] = match[0][index+1]
							}
							break
						}
					case "float":
						{
							out.Values[fieldSplit[0]] = float64(0)
							if len(match) > 0 {
								value, err := strconv.ParseFloat(match[0][index+1], 64)
								if err == nil {
									out.Values[fieldSplit[0]] = value
								}
							}
						}
					}
				}
			}
			break
		}
	}
	return out, nil
}

func (l ParsedLogEvent) ToLua() *lua.LTable {
	t := valueGoToLuaTable(l.Values)
	t.RawSetString("type", lua.LNumber(l.Type))
	t.RawSetString("raw", lua.LString(l.Raw))
	t.RawSetString("log_line", lua.LString(l.Raw))
	t.RawSetString("time", lua.LNumber(l.Time.Unix()))
	return t
}
