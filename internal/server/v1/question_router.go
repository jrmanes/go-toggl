package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jrmanes/go-toggl/pkg/question"
	"github.com/jrmanes/go-toggl/pkg/response"
)

type QuestionRouter struct {
	Repository question.Respository
}

func (qu *QuestionRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", qu.GetAllHandler)
	r.Post("/", qu.CreateHandler)
	r.Put("/{id}", qu.UpdateHandler)
	r.Delete("/{id}", qu.DeleteHandler)

	return r
}

// GetAllHandler will return an array about all questions
func (qu *QuestionRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	questions, err := qu.Repository.GetAll(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"questions": questions})
}

// CreateHandler will add a question into the database
func (qu *QuestionRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var q question.Question
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = qu.Repository.Create(ctx, &q)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusCreated, response.Map{"question": q})
}

// UpdateHandler change the values in the database using the id
func (qu *QuestionRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var q question.Question
	err = json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = qu.Repository.Update(ctx, uint(id), q)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"question": q})
}

// Delete a quetion using the id
func (qu *QuestionRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id") // in order to use chi library: https://github.com/go-chi/chi

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = qu.Repository.Delete(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"deleted": idStr})
}
