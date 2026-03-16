package handler

import (
	"log"
	"net/http"

	"github.com/davigiroux/flagkit/api/internal/cache"
	"github.com/davigiroux/flagkit/api/internal/db"
	"github.com/davigiroux/flagkit/api/internal/evaluate"
	"github.com/go-chi/chi/v5"
)

type EvalHandler struct {
	queries *db.Queries
	cache   *cache.Cache
}

func NewEvalHandler(queries *db.Queries, cache *cache.Cache) *EvalHandler {
	return &EvalHandler{queries: queries, cache: cache}
}

func (h *EvalHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	userID := r.URL.Query().Get("user_id")

	// Try cache first
	flag, err := h.cache.GetFlag(r.Context(), key)
	if err != nil {
		log.Printf("cache error for flag %q: %v", key, err)
	}
	if flag == nil {
		var err error
		flag, err = h.queries.GetFlagByKey(r.Context(), key)
		if err != nil {
			writeError(w, http.StatusNotFound, "flag not found")
			return
		}
		if err := h.cache.SetFlag(r.Context(), flag); err != nil {
			log.Printf("cache set error for flag %q: %v", key, err)
		}
	}

	result := evaluate.Evaluate(flag, userID)
	writeJSON(w, http.StatusOK, result)
}
