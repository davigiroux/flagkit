package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davigiroux/flagkit/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	pool *pgxpool.Pool
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{pool: pool}
}

func (q *Queries) ListFlags(ctx context.Context) ([]model.Flag, error) {
	rows, err := q.pool.Query(ctx, `SELECT id, key, name, description, enabled, environment, rules, created_at, updated_at FROM flags ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFlags(rows)
}

func (q *Queries) GetFlagByKey(ctx context.Context, key string) (*model.Flag, error) {
	row := q.pool.QueryRow(ctx, `SELECT id, key, name, description, enabled, environment, rules, created_at, updated_at FROM flags WHERE key = $1`, key)
	return scanFlag(row)
}

type CreateFlagInput struct {
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Enabled     bool            `json:"enabled"`
	Environment string          `json:"environment"`
	Rules       json.RawMessage `json:"rules"`
}

func (q *Queries) CreateFlag(ctx context.Context, input CreateFlagInput) (*model.Flag, error) {
	if input.Environment == "" {
		input.Environment = "development"
	}
	if input.Rules == nil {
		input.Rules = json.RawMessage("[]")
	}
	row := q.pool.QueryRow(ctx,
		`INSERT INTO flags (key, name, description, enabled, environment, rules) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, key, name, description, enabled, environment, rules, created_at, updated_at`,
		input.Key, input.Name, input.Description, input.Enabled, input.Environment, input.Rules,
	)
	return scanFlag(row)
}

type UpdateFlagInput struct {
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Enabled     *bool            `json:"enabled,omitempty"`
	Environment *string          `json:"environment,omitempty"`
	Rules       *json.RawMessage `json:"rules,omitempty"`
}

func (q *Queries) UpdateFlag(ctx context.Context, key string, input UpdateFlagInput) (*model.Flag, error) {
	existing, err := q.GetFlagByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Enabled != nil {
		existing.Enabled = *input.Enabled
	}
	if input.Environment != nil {
		existing.Environment = model.Environment(*input.Environment)
	}
	if input.Rules != nil {
		existing.Rules = *input.Rules
	}

	row := q.pool.QueryRow(ctx,
		`UPDATE flags SET name=$1, description=$2, enabled=$3, environment=$4, rules=$5, updated_at=now() WHERE key=$6 RETURNING id, key, name, description, enabled, environment, rules, created_at, updated_at`,
		existing.Name, existing.Description, existing.Enabled, existing.Environment, existing.Rules, key,
	)
	return scanFlag(row)
}

func (q *Queries) ToggleFlag(ctx context.Context, key string) (*model.Flag, error) {
	row := q.pool.QueryRow(ctx,
		`UPDATE flags SET enabled = NOT enabled, updated_at = now() WHERE key = $1 RETURNING id, key, name, description, enabled, environment, rules, created_at, updated_at`,
		key,
	)
	return scanFlag(row)
}

func (q *Queries) DeleteFlag(ctx context.Context, key string) error {
	ct, err := q.pool.Exec(ctx, `DELETE FROM flags WHERE key = $1`, key)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("flag not found")
	}
	return nil
}

// Audit logs

func (q *Queries) InsertAuditLog(ctx context.Context, flagID string, action model.AuditAction, diff json.RawMessage, actor string) error {
	_, err := q.pool.Exec(ctx,
		`INSERT INTO audit_logs (flag_id, action, diff, actor) VALUES ($1, $2, $3, $4)`,
		flagID, action, diff, actor,
	)
	return err
}

func (q *Queries) ListAuditLogs(ctx context.Context, flagID string, page, perPage int) ([]model.AuditLog, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var total int
	var countQuery string
	var countArgs []any

	if flagID != "" {
		countQuery = `SELECT COUNT(*) FROM audit_logs WHERE flag_id = $1`
		countArgs = []any{flagID}
	} else {
		countQuery = `SELECT COUNT(*) FROM audit_logs`
	}
	if err := q.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	var query string
	var args []any
	if flagID != "" {
		query = `SELECT id, flag_id, action, diff, actor, created_at FROM audit_logs WHERE flag_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []any{flagID, perPage, offset}
	} else {
		query = `SELECT id, flag_id, action, diff, actor, created_at FROM audit_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = []any{perPage, offset}
	}

	rows, err := q.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	for rows.Next() {
		var log model.AuditLog
		if err := rows.Scan(&log.ID, &log.FlagID, &log.Action, &log.Diff, &log.Actor, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}

// API keys

func (q *Queries) GetAPIKeyByHash(ctx context.Context, hash string) (string, error) {
	var id string
	err := q.pool.QueryRow(ctx, `SELECT id FROM api_keys WHERE key_hash = $1`, hash).Scan(&id)
	return id, err
}

func (q *Queries) CreateAPIKey(ctx context.Context, hash string) error {
	_, err := q.pool.Exec(ctx, `INSERT INTO api_keys (key_hash) VALUES ($1)`, hash)
	return err
}

func (q *Queries) CountAPIKeys(ctx context.Context) (int, error) {
	var count int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM api_keys`).Scan(&count)
	return count, err
}

// scan helpers

func scanFlag(row pgx.Row) (*model.Flag, error) {
	var f model.Flag
	err := row.Scan(&f.ID, &f.Key, &f.Name, &f.Description, &f.Enabled, &f.Environment, &f.Rules, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func scanFlags(rows pgx.Rows) ([]model.Flag, error) {
	var flags []model.Flag
	for rows.Next() {
		var f model.Flag
		if err := rows.Scan(&f.ID, &f.Key, &f.Name, &f.Description, &f.Enabled, &f.Environment, &f.Rules, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		flags = append(flags, f)
	}
	return flags, nil
}
