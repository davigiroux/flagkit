package model

import (
	"encoding/json"
	"time"
)

type Environment string

const (
	EnvProduction  Environment = "production"
	EnvStaging     Environment = "staging"
	EnvDevelopment Environment = "development"
)

type RuleType string

const (
	RulePercentage RuleType = "percentage"
	RuleAllowlist  RuleType = "allowlist"
)

type Rule struct {
	Type    RuleType `json:"type"`
	Value   int      `json:"value,omitempty"`
	UserIDs []string `json:"userIds,omitempty"`
}

type Flag struct {
	ID          string          `json:"id"`
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Enabled     bool            `json:"enabled"`
	Environment Environment     `json:"environment"`
	Rules       json.RawMessage `json:"rules"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func (f *Flag) ParseRules() ([]Rule, error) {
	var rules []Rule
	if len(f.Rules) == 0 {
		return rules, nil
	}
	err := json.Unmarshal(f.Rules, &rules)
	return rules, err
}

type AuditAction string

const (
	AuditCreated AuditAction = "created"
	AuditUpdated AuditAction = "updated"
	AuditDeleted AuditAction = "deleted"
	AuditToggled AuditAction = "toggled"
)

type AuditLog struct {
	ID        string          `json:"id"`
	FlagID    *string         `json:"flagId"`
	Action    AuditAction     `json:"action"`
	Diff      json.RawMessage `json:"diff"`
	Actor     string          `json:"actor"`
	CreatedAt time.Time       `json:"createdAt"`
}

type EvalResult struct {
	Enabled     bool   `json:"enabled"`
	Reason      string `json:"reason"`
	FlagKey     string `json:"flagKey"`
	EvaluatedAt string `json:"evaluatedAt"`
}
