package service

import (
	"fmt"
	"products-api/helpers"
	"products-api/models"
	"products-api/repository"
	"time"
)

type CategoryService struct {
	CategoryRepo *repository.CategoryRepository
}

// CreateCategory creates a new category
func (cs *CategoryService) CreateCategory(category *models.Category) error {
	return cs.CategoryRepo.Create(category)
}

func (cs *CategoryService) GetCategories() ([]models.Category, error) {
	cacheKey := "categories"

	var categories []models.Category

	// Try to fetch data from cache
	found, err := helpers.GetCache(cacheKey, &categories)
	if err != nil {
		return nil, err
	}

	// If cached data is found, return it
	if found {
		return categories, nil
	}

	// If cache not found, fetch data from the database
	categories, err = cs.CategoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Set cache for 10 minutes after successfully fetching from DB
	err = helpers.SetCache(cacheKey, categories, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by ID with caching
func (cs *CategoryService) GetCategoryByID(id string) (models.Category, error) {
	cacheKey := fmt.Sprintf("category:%s", id)

	var category models.Category

	// Try to fetch data from cache
	found, err := helpers.GetCache(cacheKey, &category)
	if err != nil {
		return models.Category{}, err
	}

	// If cached data is found, return it
	if found {
		return category, nil
	}

	// If cache not found, fetch data from the database
	category, err = cs.CategoryRepo.GetByID(id)
	if err != nil {
		return models.Category{}, err
	}

	// Set cache for 10 minutes after successfully fetching from DB
	err = helpers.SetCache(cacheKey, category, 10*time.Minute)
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

// UpdateProduct updates a category by ID
func (cs *CategoryService) UpdateCategory(category *models.Category) error {
	return cs.CategoryRepo.Update(category)
}

// DeleteProduct deletes a category by ID
func (cs *CategoryService) DeleteCategory(id string) error {
	return cs.CategoryRepo.Delete(id)
}
