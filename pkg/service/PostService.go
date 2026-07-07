package service

import (
	"math"
	"real-time-forum/pkg/models"
	"real-time-forum/pkg/payload/posts"
	db "real-time-forum/pkg/repositories"
)

type PostService interface {
	GetPosts(pageNumber int, pageSize int, sortBy string, sortOrder string) (posts.PostResponse, error)
	CreatePost(userID string, title string, content string, categoryId int) (posts.PostDTO, error)
}

type PostServiceImpl struct {
	db db.PostRepository
}

func NewPostService(database db.PostRepository) PostService {
	return PostServiceImpl{
		db: database,
	}
}

func mapPostToDTO(post models.Post) posts.PostDTO {
	return posts.PostDTO{
		PostId:          post.PostId,
		UserId:          post.UserId,
		Nickname:        post.Nickname,
		Title:           post.Title,
		Content:         post.Content,
		CategoryId:      post.CategoryId,
		CategoryName:    post.CategoryName,
		Score:           post.Score,
		CommentsCounter: post.CommentsCounter,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
	}
}

func (p PostServiceImpl) GetPosts(pageNumber int, pageSize int, sortBy string, sortOrder string) (posts.PostResponse, error) {
	postsModel, totalElements, err := p.db.GetPosts(pageNumber, pageSize, sortBy, sortOrder)
	if err != nil {
		return posts.PostResponse{}, err
	}

	postDTOs := make([]posts.PostDTO, len(postsModel))
	for i, post := range postsModel {
		postDTOs[i] = mapPostToDTO(post)
	}

	totalPages := int(math.Ceil(float64(totalElements) / float64(pageSize)))
	lastPage := pageNumber >= totalPages

	return posts.PostResponse{
		Posts:         postDTOs,
		PageNumber:    pageNumber,
		PageSize:      pageSize,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		LastPage:      lastPage,
	}, nil
}

func (p PostServiceImpl) CreatePost(userID string, title string, content string, categoryId int) (posts.PostDTO, error) {
	post := models.Post{
		UserId:     userID,
		Title:      title,
		Content:    content,
		CategoryId: categoryId,
	}

	createdPost, err := p.db.CreatePost(post)
	if err != nil {
		return posts.PostDTO{}, err
	}

	return mapPostToDTO(createdPost), nil
}
