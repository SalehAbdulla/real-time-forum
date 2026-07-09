package handlers

import (
	"encoding/json"
	"net/http"
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/middleware"
	"real-time-forum/pkg/models"
	"real-time-forum/pkg/payload/comment"
	"strconv"
	"strings"
)

func (re *HandlerContext) GetComments(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		re.HandleError(w, r, realtimeforum.ErrUnauthorized)
		return
	}

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

func (re *HandlerContext) CreateComments(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		re.HandleError(w, r, realtimeforum.ErrUnauthorized)
		return
	}

	if !re.parseForm(w, r) {
		return
	}

	postIdStr := strings.TrimSpace(r.FormValue("postId"))
	if postIdStr == "" {
		re.HandleError(w, r, realtimeforum.ErrMissingPostId)
		return
	}

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" || len(content) < 3 || len(content) > 100 {
		re.HandleError(w, r, realtimeforum.ErrContentLessThanThreeOrMoreThanHundard)
		return
	}

	response, err := re.CommentService.CreateComment(userID, postId, content)
	if err != nil {
		re.HandleError(w, r, err)
		return
	}

	re.App.Logger.Info("comment created successfully",
		"comment_id", response.CommentId,
		"post_id", response.PostId,
		"user_id", userID,
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.DataResponse[comment.CommentDTO]{
		Data: response,
	})
}
