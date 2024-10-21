package repository

import (
	"fmt"
	"products-api/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

// Create adds a new category to the database
func (cr *CategoryRepository) Create(category *models.Category) error {
	result := cr.DB.Create(category)
	return result.Error
}

// GetAll retrieves all categories
func (cr *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	result := cr.DB.Find(&categories)
	return categories, result.Error
}

// GetByID retrieves a category by ID
func (cr *CategoryRepository) GetByID(id string) (models.Category, error) {
	var category models.Category
	result := cr.DB.First(&category, id)
	if result.Error != nil {
		return category, fmt.Errorf("category not found")
	}
	return category, nil
}

// Update modifies an existing product
func (cr *CategoryRepository) Update(product *models.Category) error {
	result := cr.DB.Save(product)
	return result.Error
}

// Delete removes a product by ID
func (cr *CategoryRepository) Delete(id string) error {
	result := cr.DB.Delete(&models.Category{}, id)
	return result.Error
}
