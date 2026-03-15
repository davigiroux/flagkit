package audit

import (
	"encoding/json"
	"testing"

	"github.com/davigiroux/flagkit/api/internal/model"
)

func TestComputeDiff_ToggleEnabled(t *testing.T) {
	old := &model.Flag{Enabled: false, Rules: json.RawMessage("[]")}
	new_ := &model.Flag{Enabled: true, Rules: json.RawMessage("[]")}

	diff := ComputeDiff(old, new_)
	if _, ok := diff["enabled"]; !ok {
		t.Error("expected enabled in diff")
	}
	if diff["enabled"].From != false || diff["enabled"].To != true {
		t.Errorf("expected false->true, got %v->%v", diff["enabled"].From, diff["enabled"].To)
	}
}

func TestComputeDiff_NoChanges(t *testing.T) {
	flag := &model.Flag{Name: "test", Enabled: true, Rules: json.RawMessage("[]")}
	diff := ComputeDiff(flag, flag)
	if len(diff) != 0 {
		t.Errorf("expected empty diff, got %v", diff)
	}
}

func TestComputeDiff_RulesChanged(t *testing.T) {
	old := &model.Flag{Rules: json.RawMessage(`[{"type":"percentage","value":50}]`)}
	new_ := &model.Flag{Rules: json.RawMessage(`[{"type":"percentage","value":75}]`)}

	diff := ComputeDiff(old, new_)
	if _, ok := diff["rules"]; !ok {
		t.Error("expected rules in diff")
	}
}

func TestComputeDiff_MultipleFields(t *testing.T) {
	old := &model.Flag{Name: "old", Description: "old desc", Environment: "development", Rules: json.RawMessage("[]")}
	new_ := &model.Flag{Name: "new", Description: "new desc", Environment: "production", Rules: json.RawMessage("[]")}

	diff := ComputeDiff(old, new_)
	if len(diff) != 3 {
		t.Errorf("expected 3 changes, got %d: %v", len(diff), diff)
	}
}
