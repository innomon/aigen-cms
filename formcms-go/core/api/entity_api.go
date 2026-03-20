package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/formcms/formcms-go/core/services"
	"github.com/formcms/formcms-go/utils/datamodels"
)

type EntityApi struct {
	entityService services.IEntityService
}

func NewEntityApi(entityService services.IEntityService) *EntityApi {
	return &EntityApi{entityService: entityService}
}

func (a *EntityApi) Register(r chi.Router) {
	r.Route("/api/entities/{name}", func(r chi.Router) {
		r.Get("/", a.List)
		r.Post("/", a.Create)
		r.Get("/{id}", a.Get)
		r.Put("/", a.Update)
		r.Delete("/{id}", a.Delete)

		r.Route("/{id}/{attr}", func(r chi.Router) {
			r.Get("/collection", a.CollectionList)
			r.Post("/collection", a.CollectionInsert)

			r.Get("/junction", a.JunctionList)
			r.Post("/junction", a.JunctionSave)
			r.Delete("/junction", a.JunctionDelete)
		})
	})
}

func (a *EntityApi) CollectionList(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	id := chi.URLParam(r, "id")
	attr := chi.URLParam(r, "attr")

	query := r.URL.Query()
	strArgs := make(datamodels.StrArgs)
	for k, v := range query {
		strArgs[k] = v
	}

	parseResult := datamodels.ParseQuery(strArgs)
	records, total, err := a.entityService.CollectionList(r.Context(), name, id, attr, parseResult.Pagination, parseResult.Filters, parseResult.Sorts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"items": records,
		"total": total,
	}
	json.NewEncoder(w).Encode(response)
}

func (a *EntityApi) CollectionInsert(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	id := chi.URLParam(r, "id")
	attr := chi.URLParam(r, "attr")

	var record datamodels.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	savedRecord, err := a.entityService.CollectionInsert(r.Context(), name, id, attr, record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(savedRecord)
}

func (a *EntityApi) JunctionList(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	id := chi.URLParam(r, "id")
	attr := chi.URLParam(r, "attr")
	exclude := r.URL.Query().Get("exclude") == "true"

	query := r.URL.Query()
	strArgs := make(datamodels.StrArgs)
	for k, v := range query {
		strArgs[k] = v
	}

	parseResult := datamodels.ParseQuery(strArgs)
	records, total, err := a.entityService.JunctionList(r.Context(), name, id, attr, exclude, parseResult.Pagination, parseResult.Filters, parseResult.Sorts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"items": records,
		"total": total,
	}
	json.NewEncoder(w).Encode(response)
}

func (a *EntityApi) JunctionSave(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	id := chi.URLParam(r, "id")
	attr := chi.URLParam(r, "attr")

	var targetIds []interface{}
	if err := json.NewDecoder(r.Body).Decode(&targetIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.entityService.JunctionSave(r.Context(), name, id, attr, targetIds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *EntityApi) JunctionDelete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	id := chi.URLParam(r, "id")
	attr := chi.URLParam(r, "attr")

	var targetIds []interface{}
	if err := json.NewDecoder(r.Body).Decode(&targetIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.entityService.JunctionDelete(r.Context(), name, id, attr, targetIds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *EntityApi) List(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	query := r.URL.Query()
	strArgs := make(datamodels.StrArgs)
	for k, v := range query {
		strArgs[k] = v
	}

	parseResult := datamodels.ParseQuery(strArgs)
	records, total, err := a.entityService.List(r.Context(), name, parseResult.Pagination, parseResult.Filters, parseResult.Sorts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"items": records,
		"total": total,
	}
	json.NewEncoder(w).Encode(response)
}

func (a *EntityApi) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	id := chi.URLParam(r, "id")

	record, err := a.entityService.Single(r.Context(), name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(record)
}

func (a *EntityApi) Create(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var record datamodels.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	savedRecord, err := a.entityService.Insert(r.Context(), name, record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(savedRecord)
}

func (a *EntityApi) Update(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var record datamodels.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedRecord, err := a.entityService.Update(r.Context(), name, record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updatedRecord)
}

func (a *EntityApi) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	id := chi.URLParam(r, "id")

	if err := a.entityService.Delete(r.Context(), name, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
