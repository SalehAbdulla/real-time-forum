package repositories

import (
	"real-time-forum/pkg/payload/posts"
)

type PostRepository interface {
	Getposts(pageNumber int, Integer int, sortBy string, sortOrder string) (posts.PostResponse, error)
}

func (db *DB) Getposts(pageNumber int, Integer int, sortBy string, sortOrder string) (posts.PostResponse, error) {
	return posts.PostResponse{}, nil
}
