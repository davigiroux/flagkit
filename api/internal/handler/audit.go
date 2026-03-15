package handler

import (
	"net/http"
	"strconv"

	"github.com/davigiroux/flagkit/api/internal/db"
)

type AuditHandler struct {
	queries *db.Queries
}

func NewAuditHandler(queries *db.Queries) *AuditHandler {
	return &AuditHandler{queries: queries}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	flagID := r.URL.Query().Get("flag_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	logs, total, err := h.queries.ListAuditLogs(r.Context(), flagID, page, perPage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list audit logs")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data":  logs,
		"total": total,
		"page":  page,
	})
}
