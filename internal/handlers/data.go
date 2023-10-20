package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"dev-challenge/internal/models"
	"dev-challenge/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type ExcelLikeHandler struct {
	ELS services.ExcelLikeService
	Log logrus.FieldLogger
}

func (h *ExcelLikeHandler) RegisterRoutes(router chi.Router) {
	router.Post("/{sheet_id}/{cell_id}", h.addValue)
	router.Get("/{sheet_id}/{cell_id}", h.getValue)
	router.Get("/{sheet_id}", h.getAllValues)
}

func (h *ExcelLikeHandler) getValue(w http.ResponseWriter, r *http.Request) {
	sheetID := chi.URLParam(r, "sheet_id")
	cellID := chi.URLParam(r, "cell_id")
	if !containsOnlyURLAllowedChars(strings.ToLower(sheetID)) || !containsOnlyURLAllowedChars(strings.ToLower(cellID)) {
		h.Log.Error("not correct data in params")
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, models.Error("not correct params", http.StatusNotFound))
		return
	}
	cellInput, err := h.ELS.GetCellInput(r.Context(), strings.ToLower(sheetID), strings.ToLower(cellID))
	if err != nil {
		h.Log.WithError(err)
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, models.Error("store not responded", http.StatusNotFound))
		return
	}
	if cellInput.Value == "" {
		h.Log.Error(fmt.Sprintf("value on sheetID=%s and cellID=%s not found", sheetID, cellID))
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, models.Error("value not found", http.StatusNotFound))
		return
	}
	render.JSON(w, r, cellInput)
}

func (h *ExcelLikeHandler) addValue(w http.ResponseWriter, r *http.Request) {
	sheetID := chi.URLParam(r, "sheet_id")
	cellID := chi.URLParam(r, "cell_id")
	if !containsOnlyURLAllowedChars(strings.ToLower(sheetID)) || !containsOnlyURLAllowedChars(strings.ToLower(cellID)) {
		h.Log.Error("not correct data in params")
		w.WriteHeader(http.StatusUnprocessableEntity)
		render.JSON(w, r, models.Error("not correct params", http.StatusUnprocessableEntity))
		return
	}
	var requestBody *models.Data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Log.Errorf(fmt.Sprintf("not correct request body: %s", requestBody))
		w.WriteHeader(http.StatusUnprocessableEntity)
		render.JSON(w, r, models.Error("can't read request body", http.StatusUnprocessableEntity))
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		render.JSON(w, r, models.Error("can't unmarshal request body", http.StatusUnprocessableEntity))
		return
	}

	resp, err := h.ELS.AddCellInputTX(r.Context(), strings.ToLower(sheetID), strings.ToLower(cellID), requestBody)
	if err != nil {
		h.Log.WithError(err).Error("failed to add value")
		w.WriteHeader(http.StatusUnprocessableEntity)
		render.JSON(w, r, models.ErrorPOSTResponse(requestBody.Value))
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resp)
}

func (h *ExcelLikeHandler) getAllValues(w http.ResponseWriter, r *http.Request) {
	sheetID := chi.URLParam(r, "sheet_id")
	if !containsOnlyURLAllowedChars(strings.ToLower(sheetID)) {
		h.Log.Error("not correct data in params")
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, models.Error("not correct params", http.StatusNotFound))
		return
	}
	cellInput, err := h.ELS.GetSheetInput(r.Context(), strings.ToLower(sheetID))
	if err != nil {
		h.Log.WithError(err)
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, models.Error("store not responded", http.StatusNotFound))
		return
	}
	if cellInput == nil {
		h.Log.Error(fmt.Sprintf("sheetID=%s not found", sheetID))
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, models.Error("value not found", http.StatusNotFound))
		return
	}
	render.JSON(w, r, cellInput)
}

func containsOnlyURLAllowedChars(s string) bool {
	pattern := "^[a-z0-9-_.~%!$&'()*+,;=:@/\\[\\]?#]+$"
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}
	return matched
}
