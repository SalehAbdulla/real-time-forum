package repositories

import (
	"real-time-forum/pkg/models"
)

type CommentRepository interface {
	GetComments(postId int, pageNumber int, pageSize int, sortBy string, sortOrder string) ([]models.Comment, error)
}

func (db *DB) GetComments(postId int, pageNumber int, pageSize int, sortBy string, sortOrder string) ([]models.Comment, error) {

	return []models.Comment{}, nil
}
