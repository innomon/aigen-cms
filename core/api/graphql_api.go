package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/formcms/formcms-go/core/services"
)

type GraphQLApi struct {
	graphqlService services.IGraphQLService
}

func NewGraphQLApi(graphqlService services.IGraphQLService) *GraphQLApi {
	return &GraphQLApi{graphqlService: graphqlService}
}

func (a *GraphQLApi) Register(r chi.Router) {
	r.Post("/api/graphql", a.Query)
}

type graphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (a *GraphQLApi) Query(w http.ResponseWriter, r *http.Request) {
	var req graphqlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := a.graphqlService.Query(r.Context(), req.Query, req.Variables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
