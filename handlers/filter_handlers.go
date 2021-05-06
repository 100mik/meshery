package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
)

// FilterFileRequestHandler will handle requests of both type GET and POST
// on the route /api/experimental/filter
func (h *Handler) FilterFileRequestHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	if r.Method == http.MethodGet {
		h.GetMesheryFiltersHandler(rw, r, prefObj, user, provider)
		return
	}

	if r.Method == http.MethodPost {
		h.SaveFilterFile(rw, r, prefObj, user, provider)
		return
	}
}

// GetMesheryFilterHandler will handle requests type GET
// on the route /api/experimental/filter/{id}
func (h *Handler) GetMesheryFilterHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	filterID := mux.Vars(r)["id"]

	resp, err := provider.GetMesheryFilter(r, filterID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to get the filter: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// DeleteMesheryFilterHandler deletes a filter with the given id
func (h *Handler) DeleteMesheryFilterHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	filterID := mux.Vars(r)["id"]

	resp, err := provider.DeleteMesheryFilter(r, filterID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to delete the filter: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// ImportFilterFileGithub will import filterfile file from github
func (h *Handler) ImportFilterFileGithub(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	owner := mux.Vars(r)["owner"]
	repo := mux.Vars(r)["repo"]
	path := mux.Vars(r)["path"]

	cont, err := provider.ImportFilterFileGithub(r, owner, repo, path)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", err)
		return
	}

	fmt.Fprintf(rw, "%s", cont)
}

func (h *Handler) SaveFilterFile(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var parsedBody *models.MesheryFilter

	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, "failed to get user token", http.StatusInternalServerError)
		return
	}

	resp, err := provider.SaveMesheryFilter(token, parsedBody)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to save the filter: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// GetMesheryFiltersHandler returns the list of all the patterns saved by the current user
func (h *Handler) GetMesheryFiltersHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()

	resp, err := provider.GetMesheryFilters(r, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to fetch the filters: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}
