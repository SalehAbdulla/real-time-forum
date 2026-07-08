package service

import (
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/payload/reaction"
	db "real-time-forum/pkg/repositories"
)

type ReactionService interface {
	UpsertReaction(userId string, entityType string, entityId int, score int) (reaction.ReactionResponse, error)
}

type ReactionServiceImpl struct {
	db db.ReactionRepository
}

func NewReactionService(database db.ReactionRepository) ReactionService {
	return ReactionServiceImpl{
		db: database,
	}
}

func (r ReactionServiceImpl) UpsertReaction(userId string, entityType string, entityId int, score int) (reaction.ReactionResponse, error) {
	if entityType != "post" && entityType != "comment" {
		return reaction.ReactionResponse{}, realtimeforum.ErrBadRequest
	}

	if score != 1 && score != -1 {
		return reaction.ReactionResponse{}, realtimeforum.ErrBadRequest
	}

	if entityId < 1 {
		return reaction.ReactionResponse{}, realtimeforum.ErrBadRequest
	}

	totalScore, err := r.db.UpsertReaction(userId, entityType, entityId, score)
	if err != nil {
		return reaction.ReactionResponse{}, err
	}

	return reaction.ReactionResponse{
		TotalScore: totalScore,
	}, nil
}