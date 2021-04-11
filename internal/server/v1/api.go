package v1

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jrmanes/go-toggl/internal/data"
)

// New returns the API V1 Handler with configuration.
func New() http.Handler {
	r := chi.NewRouter()

	// Mount quetions routes
	qu := &QuestionRouter{
		Repository: &data.QuestionRepository{
			Data: data.New(),
		},
	}
	r.Mount("/q", qu.Routes())

	return r
}
