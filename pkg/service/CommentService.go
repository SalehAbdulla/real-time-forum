package service

import (
	"math"
	"real-time-forum/pkg/payload/comment"
	db "real-time-forum/pkg/repositories"
)

type CommentService interface {
	GetComments(postId int, pageNumber int, pageSize int, sortBy string, sortOrder string) (comment.CommentResponse, error)
}

type CommentServiceImpl struct {
	db db.CommentRepository
}

func NewCommentService(database db.CommentRepository) CommentService {
	return CommentServiceImpl{
		db: database,
	}
}

func (c CommentServiceImpl) GetComments(postId int, pageNumber int, pageSize int, sortBy string, sortOrder string) (comment.CommentResponse, error) {
	comments, totalElements, err := c.db.GetComments(postId, pageNumber, pageSize, sortBy, sortOrder)
	if err != nil {
		return comment.CommentResponse{}, err
	}

	dtos := make([]comment.CommentDTO, len(comments))
	for i, com := range comments {
		dtos[i] = comment.CommentDTO{
			CommentId:   com.CommentId,
			PostId:      com.PostId,
			UserId:      com.UserId,
			CommentText: com.CommentText,
			CreatedAt:   com.CreatedAt,
		}
	}

	totalPages := int(math.Ceil(float64(totalElements) / float64(pageSize)))
	lastPage := pageNumber >= totalPages

	return comment.CommentResponse{
		Comments:      dtos,
		PageNumber:    pageNumber,
		PageSize:      pageSize,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		LastPage:      lastPage,
	}, nil
}