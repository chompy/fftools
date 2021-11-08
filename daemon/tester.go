package main

import (
	"bufio"
	"io"
	"sort"
	"strings"
	"time"
)

type testerLine struct {
	time    time.Time
	logLine string
}

func testerParse(data io.Reader) []LogLine {
	startTime := time.Time{}
	scanner := bufio.NewScanner(data)
	lines := make([]testerLine, 0)
	for scanner.Scan() {
		logLine := scanner.Text()
		split := strings.Split(logLine, " ")
		timestampStr := strings.Trim(split[0], "[]")
		timestamp, err := time.Parse("15:04:05.999999999", timestampStr)
		if err != nil {
			logWarn(err.Error())
			continue
		}
		if startTime.IsZero() || timestamp.Before(startTime) {
			startTime = timestamp
		}
		lines = append(lines, testerLine{
			time:    timestamp,
			logLine: logLine,
		})
	}
	timeDiff := time.Since(startTime)
	encounterID := uint32(time.Now().Unix() / 4)
	out := make([]LogLine, 0)
	for _, line := range lines {
		out = append(out, LogLine{
			EncounterID: encounterID,
			Time:        line.time.Add(timeDiff),
			LogLine:     line.logLine,
		})
	}
	sort.Slice(out, func(p, q int) bool {
		return out[p].Time.Before(out[q].Time)
	})
	return out
}

func testerReplay(logLines []LogLine) {
	logInfo("Replay %d log lines.", len(logLines))
	// generate encounter
	eventListenerDispatch(
		"act:encounter",
		Encounter{
			ID:        logLines[0].EncounterID,
			StartTime: time.Now(),
			Active:    true,
			Zone:      "Test Area",
		},
	)
	// read combatants
	combatants := make([]Combatant, 0)
	for _, logLine := range logLines {
		pll, err := ParseLogEvent(logLine)
		if err != nil {
			logWarn(err.Error())
			continue
		}
		switch pll.Type {
		case LogTypeNetworkAbility, LogTypeNetworkAOEAbility:
			{
				combatant := Combatant{
					ID:          int32(pll.Values["source_id"].(int)),
					Name:        pll.Values["source_name"].(string),
					EncounterID: logLines[0].EncounterID,
					Job:         "War",
				}
				hasCombatant := false
				for i := range combatants {
					if combatants[i].ID == combatant.ID {
						hasCombatant = true
						break
					}
				}
				if !hasCombatant {
					combatants = append(combatants, combatant)
				}
				break
			}
		}
	}
	for _, combatant := range combatants {
		eventListenerDispatch("act:combatant", combatant)
	}
	// send log lines in real time
	time.Sleep(time.Second)
	for i := range logLines {
		pll, err := ParseLogEvent(logLines[i])
		if err != nil {
			logWarn(err.Error())
			continue
		}
		eventListenerDispatch("act:log_line", pll)
		if len(logLines) > i+1 {
			time.Sleep(logLines[i+1].Time.Sub(logLines[i].Time))
		}
	}
}
