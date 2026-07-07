package handlers

import (
	"encoding/json"
	"net/http"
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/middleware"
	"real-time-forum/pkg/models"
	"real-time-forum/pkg/payload/posts"
	"strconv"
	"strings"
)

func (re *Repository) GetPosts(w http.ResponseWriter, r *http.Request) {
	pageNumberStr := r.URL.Query().Get("page")
	if pageNumberStr == "" {
		pageNumberStr = "1"
	}

	pageSizeStr := r.URL.Query().Get("size")
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}

	sortBy := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("sortBy")))
	if sortBy == "" {
		sortBy = "createdAt"
	}

	sortOrder := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("sortOrder")))
	if sortOrder == "" {
		sortOrder = "desc"
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber < 1 {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	validSortBy := map[string]bool{
		"createdat": true,
		"title":     true,
		"score":     true,
	}

	if !validSortBy[sortBy] {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	validSortOrder := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validSortOrder[sortOrder] {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

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

func (re *Repository) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		re.HandleError(w, r, realtimeforum.ErrUnauthorized)
		return
	}

	if !re.parseForm(w, r) {
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categoryIdStr := strings.TrimSpace(r.FormValue("categoryId"))

	if title == "" || len(title) > 30 {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	if content == "" || len(content) > 100  {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	if categoryIdStr == "" {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil || categoryId < 1 {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	response, err := re.PostService.CreatePost(userID, title, content, categoryId)
	if err != nil {
		re.HandleError(w, r, err)
		return
	}

	re.App.Logger.Info("post created successfully",
		"post_id", response.PostId,
		"user_id", userID,
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.DataResponse[posts.PostDTO]{
		Data: response,
	})
}
