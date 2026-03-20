package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/formcms/formcms-go/core/services"
	"github.com/formcms/formcms-go/utils/datamodels"
)

type PageApi struct {
	pageService services.IPageService
}

func NewPageApi(pageService services.IPageService) *PageApi {
	return &PageApi{pageService: pageService}
}

func (a *PageApi) Register(r chi.Router) {
	// Catch all for pages
	r.Get("/*", a.Render)
}

func (a *PageApi) Render(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	query := r.URL.Query()
	strArgs := make(datamodels.StrArgs)
	for k, v := range query {
		strArgs[k] = v
	}

	html, err := a.pageService.Render(r.Context(), path, strArgs)
	if err != nil {
		if path == "/" || path == "" {
			http.Redirect(w, r, "/admin/list.html", http.StatusFound)
			return
		}
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
