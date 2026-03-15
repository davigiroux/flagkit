package audit

import (
	"encoding/json"

	"github.com/davigiroux/flagkit/api/internal/model"
)

type FieldChange struct {
	From any `json:"from"`
	To   any `json:"to"`
}

func ComputeDiff(old, new_ *model.Flag) map[string]FieldChange {
	diff := make(map[string]FieldChange)

	if old.Name != new_.Name {
		diff["name"] = FieldChange{From: old.Name, To: new_.Name}
	}
	if old.Description != new_.Description {
		diff["description"] = FieldChange{From: old.Description, To: new_.Description}
	}
	if old.Enabled != new_.Enabled {
		diff["enabled"] = FieldChange{From: old.Enabled, To: new_.Enabled}
	}
	if old.Environment != new_.Environment {
		diff["environment"] = FieldChange{From: old.Environment, To: new_.Environment}
	}
	if string(old.Rules) != string(new_.Rules) {
		var oldRules, newRules any
		json.Unmarshal(old.Rules, &oldRules)
		json.Unmarshal(new_.Rules, &newRules)
		diff["rules"] = FieldChange{From: oldRules, To: newRules}
	}

	return diff
}

func MarshalDiff(diff map[string]FieldChange) json.RawMessage {
	data, _ := json.Marshal(diff)
	return data
}

func FlagSnapshot(flag *model.Flag) json.RawMessage {
	var rules any
	json.Unmarshal(flag.Rules, &rules)

	diff := map[string]FieldChange{
		"key":         {From: nil, To: flag.Key},
		"name":        {From: nil, To: flag.Name},
		"description": {From: nil, To: flag.Description},
		"enabled":     {From: nil, To: flag.Enabled},
		"environment": {From: nil, To: string(flag.Environment)},
		"rules":       {From: nil, To: rules},
	}
	return MarshalDiff(diff)
}

func FlagDeleteSnapshot(flag *model.Flag) json.RawMessage {
	var rules any
	json.Unmarshal(flag.Rules, &rules)

	diff := map[string]FieldChange{
		"key":         {From: flag.Key, To: nil},
		"name":        {From: flag.Name, To: nil},
		"description": {From: flag.Description, To: nil},
		"enabled":     {From: flag.Enabled, To: nil},
		"environment": {From: string(flag.Environment), To: nil},
		"rules":       {From: rules, To: nil},
	}
	return MarshalDiff(diff)
}
