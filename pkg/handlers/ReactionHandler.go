package handlers

import (
	"encoding/json"
	"net/http"
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/middleware"
	"real-time-forum/pkg/models"
	"real-time-forum/pkg/payload/reaction"
	"strconv"
	"strings"
)

type ReactRequest struct {
	EntityType string `json:"entityType"`
	EntityId   int    `json:"entityId"`
	Score      int    `json:"score"`
}

func (re *Repository) React(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		re.HandleError(w, r, realtimeforum.ErrUnauthorized)
		return
	}
	req := ReactRequest{}
	req.EntityType = strings.TrimSpace(strings.ToLower(r.FormValue("entityType")))
	entityId, err := strconv.Atoi(r.FormValue("entityId"))
	if err != nil {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}
	req.EntityId = entityId
	score, err := strconv.Atoi(r.FormValue("score"))
	if err != nil {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}
	req.Score = score

	response, err := re.ReactService.UpsertReaction(userID, req.EntityType, req.EntityId, req.Score)
	if err != nil {
		re.HandleError(w, r, err)
		return
	}

	re.App.Logger.Info("reaction upserted successfully",
		"user_id", userID,
		"entity_type", req.EntityType,
		"entity_id", req.EntityId,
		"total_score", response.TotalScore,
	)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.DataResponse[reaction.ReactionResponse]{
		Data: response,
	})
}
