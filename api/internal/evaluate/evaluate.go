package evaluate

import (
	"time"

	"github.com/davigiroux/flagkit/api/internal/hash"
	"github.com/davigiroux/flagkit/api/internal/model"
)

func Evaluate(flag *model.Flag, userID string) model.EvalResult {
	result := model.EvalResult{
		FlagKey:     flag.Key,
		EvaluatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if !flag.Enabled {
		result.Reason = "flag_disabled"
		return result
	}

	rules, err := flag.ParseRules()
	if err != nil {
		result.Reason = "invalid_rules"
		return result
	}

	for _, rule := range rules {
		switch rule.Type {
		case model.RuleAllowlist:
			for _, id := range rule.UserIDs {
				if id == userID {
					result.Enabled = true
					result.Reason = "allowlist"
					return result
				}
			}
		case model.RulePercentage:
			bucket := hash.Consistent(flag.Key, userID)
			if bucket < rule.Value {
				result.Enabled = true
				result.Reason = "rollout"
				return result
			}
		}
	}

	result.Reason = "no_match"
	return result
}
