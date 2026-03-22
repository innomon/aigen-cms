package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/core/services"
)

type EngagementApi struct {
	engagementService services.IEngagementService
}

func NewEngagementApi(engagementService services.IEngagementService) *EngagementApi {
	return &EngagementApi{engagementService: engagementService}
}

func (a *EngagementApi) Register(r chi.Router) {
	r.Post("/api/engagements/track", a.Track)
}

func (a *EngagementApi) Track(w http.ResponseWriter, r *http.Request) {
	var status descriptors.EngagementStatus
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In a real app, get userId from auth or cookie
	if status.UserId == "" {
		status.UserId = "anonymous"
	}

	if err := a.engagementService.Track(r.Context(), &status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
