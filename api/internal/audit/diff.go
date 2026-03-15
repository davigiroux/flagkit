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
	data, _ := json.Marshal(flag)
	return data
}
