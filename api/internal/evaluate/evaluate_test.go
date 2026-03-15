package evaluate

import (
	"encoding/json"
	"testing"

	"github.com/davigiroux/flagkit/api/internal/model"
)

func makeFlag(key string, enabled bool, rules []model.Rule) *model.Flag {
	rulesJSON, _ := json.Marshal(rules)
	return &model.Flag{
		ID:      "test-id",
		Key:     key,
		Name:    "Test Flag",
		Enabled: enabled,
		Rules:   rulesJSON,
	}
}

func TestEvaluate_DisabledFlag(t *testing.T) {
	flag := makeFlag("test", false, nil)
	result := Evaluate(flag, "user-1")
	if result.Enabled {
		t.Error("expected disabled")
	}
	if result.Reason != "flag_disabled" {
		t.Errorf("expected reason flag_disabled, got %s", result.Reason)
	}
}

func TestEvaluate_AllowlistMatch(t *testing.T) {
	flag := makeFlag("test", true, []model.Rule{
		{Type: model.RuleAllowlist, UserIDs: []string{"user-1", "user-2"}},
	})
	result := Evaluate(flag, "user-1")
	if !result.Enabled {
		t.Error("expected enabled")
	}
	if result.Reason != "allowlist" {
		t.Errorf("expected reason allowlist, got %s", result.Reason)
	}
}

func TestEvaluate_AllowlistNoMatch(t *testing.T) {
	flag := makeFlag("test", true, []model.Rule{
		{Type: model.RuleAllowlist, UserIDs: []string{"user-1"}},
	})
	result := Evaluate(flag, "user-99")
	if result.Enabled {
		t.Error("expected disabled")
	}
	if result.Reason != "no_match" {
		t.Errorf("expected reason no_match, got %s", result.Reason)
	}
}

func TestEvaluate_Percentage100(t *testing.T) {
	flag := makeFlag("test", true, []model.Rule{
		{Type: model.RulePercentage, Value: 100},
	})
	result := Evaluate(flag, "any-user")
	if !result.Enabled {
		t.Error("expected enabled at 100%")
	}
	if result.Reason != "rollout" {
		t.Errorf("expected reason rollout, got %s", result.Reason)
	}
}

func TestEvaluate_Percentage0(t *testing.T) {
	flag := makeFlag("test", true, []model.Rule{
		{Type: model.RulePercentage, Value: 0},
	})
	result := Evaluate(flag, "any-user")
	if result.Enabled {
		t.Error("expected disabled at 0%")
	}
}

func TestEvaluate_PercentageDeterministic(t *testing.T) {
	flag := makeFlag("my-flag", true, []model.Rule{
		{Type: model.RulePercentage, Value: 50},
	})
	r1 := Evaluate(flag, "user-abc")
	r2 := Evaluate(flag, "user-abc")
	if r1.Enabled != r2.Enabled {
		t.Error("expected deterministic result")
	}
}

func TestEvaluate_RuleOrder_AllowlistBeforePercentage(t *testing.T) {
	flag := makeFlag("test", true, []model.Rule{
		{Type: model.RuleAllowlist, UserIDs: []string{"vip-user"}},
		{Type: model.RulePercentage, Value: 0},
	})
	result := Evaluate(flag, "vip-user")
	if !result.Enabled {
		t.Error("allowlist should match before percentage")
	}
	if result.Reason != "allowlist" {
		t.Errorf("expected allowlist, got %s", result.Reason)
	}
}

func TestEvaluate_NoRules(t *testing.T) {
	flag := makeFlag("test", true, nil)
	result := Evaluate(flag, "user-1")
	if result.Enabled {
		t.Error("expected disabled with no rules")
	}
	if result.Reason != "no_match" {
		t.Errorf("expected no_match, got %s", result.Reason)
	}
}

func TestEvaluate_FlagKeyInResult(t *testing.T) {
	flag := makeFlag("my-flag", false, nil)
	result := Evaluate(flag, "user-1")
	if result.FlagKey != "my-flag" {
		t.Errorf("expected flagKey my-flag, got %s", result.FlagKey)
	}
}
