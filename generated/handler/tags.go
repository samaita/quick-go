package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samaita/quick-go/generated/model"
	"github.com/samaita/quick-go/generated/repo"
)

type TagRepo interface {
	List(ctx context.Context, limit, offset int) ([]*model.Tag, error)
	GetByID(ctx context.Context, id int64) (*model.Tag, error)
	Create(ctx context.Context, m *model.Tag) (int64, error)
	Update(ctx context.Context, id int64, m *model.Tag) error
	Delete(ctx context.Context, id int64) error
}

type TagHandler struct {
	repo TagRepo
}

func NewTagHandler(r *repo.TagRepo) *TagHandler {
	return &TagHandler{repo: r}
}

func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 20)
	offset := queryInt(r, "offset", 0)

	items, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []*model.Tag{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *TagHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDInt64(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.repo.GetByID(r.Context(), int64(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var m model.Tag
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.repo.Create(r.Context(), &m)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDInt64(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var m model.Tag
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(r.Context(), int64(id), &m); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDInt64(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(r.Context(), int64(id)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
