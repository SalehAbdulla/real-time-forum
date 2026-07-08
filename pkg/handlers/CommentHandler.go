package handlers

import (
	"encoding/json"
	"net/http"
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/models"
	"real-time-forum/pkg/payload/comment"
	"strconv"
	"strings"
)

func (re *Repository) GetComments(w http.ResponseWriter, r *http.Request) {
	postId := strings.TrimSpace(r.URL.Query().Get("postId"))
	if postId == "" {
		re.HandleError(w, r, realtimeforum.ErrMissingPostId)
		return
	}

	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	pageNumberStr := r.URL.Query().Get("page")
	if pageNumberStr == "" {
		pageNumberStr = "1"
	}

	pageSizeStr := r.URL.Query().Get("size")
	if pageSizeStr == "" {
		pageSizeStr = "10"
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

	response, err := re.CommentService.GetComments(postIdInt, pageNumber, pageSize, "createdAt", "desc")
	if err != nil {
		re.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.DataResponse[comment.CommentResponse]{
		Data: response,
	})
}
