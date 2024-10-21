package controllers

import (
	"net/http"

	"products-api/models"
	"products-api/service"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	Service *service.CategoryService
}

// CreateCategory handles the creation of a new category
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Service.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategories retrieves a list of categories
func (cc *CategoryController) GetCategories(c *gin.Context) {
	categories, err := cc.Service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID retrieves a category by ID
func (cc *CategoryController) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := cc.Service.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateProduct updates a category by ID
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Service.UpdateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category by ID
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := cc.Service.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
