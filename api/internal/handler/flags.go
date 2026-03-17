package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/davigiroux/flagkit/api/internal/audit"
	"github.com/davigiroux/flagkit/api/internal/cache"
	"github.com/davigiroux/flagkit/api/internal/db"
	"github.com/davigiroux/flagkit/api/internal/middleware"
	"github.com/davigiroux/flagkit/api/internal/model"
	"github.com/go-chi/chi/v5"
)

type FlagHandler struct {
	queries *db.Queries
	cache   *cache.Cache
}

func NewFlagHandler(queries *db.Queries, cache *cache.Cache) *FlagHandler {
	return &FlagHandler{queries: queries, cache: cache}
}

func (h *FlagHandler) List(w http.ResponseWriter, r *http.Request) {
	flags, err := h.queries.ListFlags(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list flags")
		return
	}
	if flags == nil {
		writeJSON(w, http.StatusOK, []any{})
		return
	}
	writeJSON(w, http.StatusOK, flags)
}

func (h *FlagHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	flag, err := h.queries.GetFlagByKey(r.Context(), key)
	if err != nil {
		writeError(w, http.StatusNotFound, "flag not found")
		return
	}
	writeJSON(w, http.StatusOK, flag)
}

func (h *FlagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input db.CreateFlagInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Key == "" || input.Name == "" {
		writeError(w, http.StatusBadRequest, "key and name are required")
		return
	}
	flag, err := h.queries.CreateFlag(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusConflict, "flag already exists or invalid data")
		return
	}

	if err := h.queries.InsertAuditLog(r.Context(), flag.ID, model.AuditCreated, audit.FlagSnapshot(flag), middleware.GetActor(r.Context())); err != nil {
		log.Printf("audit log error (create flag %q): %v", flag.Key, err)
	}

	writeJSON(w, http.StatusCreated, flag)
}

func (h *FlagHandler) Update(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	oldFlag, err := h.queries.GetFlagByKey(r.Context(), key)
	if err != nil {
		writeError(w, http.StatusNotFound, "flag not found")
		return
	}

	var input db.UpdateFlagInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	flag, err := h.queries.UpdateFlag(r.Context(), key, input)
	if err != nil {
		writeError(w, http.StatusNotFound, "flag not found")
		return
	}

	diff := audit.ComputeDiff(oldFlag, flag)
	if err := h.queries.InsertAuditLog(r.Context(), flag.ID, model.AuditUpdated, audit.MarshalDiff(diff), middleware.GetActor(r.Context())); err != nil {
		log.Printf("audit log error (update flag %q): %v", key, err)
	}

	if err := h.cache.InvalidateFlag(r.Context(), key); err != nil {
		log.Printf("cache invalidation error for flag %q: %v", key, err)
	}
	writeJSON(w, http.StatusOK, flag)
}

func (h *FlagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	flag, err := h.queries.GetFlagByKey(r.Context(), key)
	if err != nil {
		writeError(w, http.StatusNotFound, "flag not found")
		return
	}

	if err := h.queries.DeleteFlag(r.Context(), key); err != nil {
		writeError(w, http.StatusNotFound, "flag not found")
		return
	}

	if err := h.queries.InsertAuditLog(r.Context(), flag.ID, model.AuditDeleted, audit.FlagDeleteSnapshot(flag), middleware.GetActor(r.Context())); err != nil {
		log.Printf("audit log error (delete flag %q): %v", key, err)
	}

	if err := h.cache.InvalidateFlag(r.Context(), key); err != nil {
		log.Printf("cache invalidation error for flag %q: %v", key, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FlagHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	oldFlag, err := h.queries.GetFlagByKey(r.Context(), key)
	if err != nil {
		writeError(w, http.StatusNotFound, "flag not found")
		return
	}

	flag, err := h.queries.ToggleFlag(r.Context(), key)
	if err != nil {
		writeError(w, http.StatusNotFound, "flag not found")
		return
	}

	diff := audit.ComputeDiff(oldFlag, flag)
	if err := h.queries.InsertAuditLog(r.Context(), flag.ID, model.AuditToggled, audit.MarshalDiff(diff), middleware.GetActor(r.Context())); err != nil {
		log.Printf("audit log error (toggle flag %q): %v", key, err)
	}

	if err := h.cache.InvalidateFlag(r.Context(), key); err != nil {
		log.Printf("cache invalidation error for flag %q: %v", key, err)
	}
	writeJSON(w, http.StatusOK, flag)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
