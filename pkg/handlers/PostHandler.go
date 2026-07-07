package handlers

import (
	"encoding/json"
	"net/http"
	"real-time-forum/pkg/models"
	"real-time-forum/pkg/payload/posts"
	"strconv"
	"strings"
)

func (re *Repository) GetPosts(w http.ResponseWriter, r *http.Request) {

	pageNumber, err := strconv.Atoi(r.PathValue("pageNumber"))
	if err != nil {
		re.HandleError(w, r, err)
	}

	pageSize, err := strconv.Atoi(r.PathValue("pageSize"))
	if err != nil {
		re.HandleError(w, r, err)
	}

	sortBy := strings.TrimSpace(strings.ToLower(r.PathValue("sortBy")))
	if err != nil {
		re.HandleError(w, r, err)
	}

	sortOrder := strings.TrimSpace(strings.ToLower(r.PathValue("sortOrder")))

	response, err := re.PostService.GetPosts(pageNumber, pageSize, sortBy, sortOrder)

	if err != nil {
		re.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.DataResponse[posts.PostResponse]{
		Data: response,
	})
}
