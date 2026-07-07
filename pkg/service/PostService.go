package service

import (
	"real-time-forum/pkg/payload/posts"
	db "real-time-forum/pkg/repositories"
)

type PostService interface {
	GetPosts(pageNumber int, Integer int, sortBy string, sortOrder string) (posts.PostResponse, error)
}

type PostServiceImpl struct {
	db db.PostRepository
}

func NewPostService(database db.PostRepository) PostService {
	return PostServiceImpl{
		db: database,
	}
}

func (p PostServiceImpl) GetPosts(pageNumber int, Integer int, sortBy string, sortOrder string) (posts.PostResponse, error) {

	// category, err := c.db.GetPosts()

	return posts.PostResponse{}, nil
}
