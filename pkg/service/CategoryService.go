package service

import (
	"real-time-forum/pkg/models"
	db "real-time-forum/pkg/repositories"
)

type CategoryService interface {
	GetCategories() ([]models.Category, error)
}

type CategoryServiceImpl struct {
	db db.CategoryRepository
}

func NewCategoryService(database db.CategoryRepository) CategoryService {
	return CategoryServiceImpl{
		db: database,
	}
}

func (c CategoryServiceImpl) GetCategories() ([]models.Category, error) {
	return c.db.GetCategories()
}
