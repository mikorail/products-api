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

// CreateCategory godoc
// @Summary      Create Category
// @Description  Create a new category
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category  body      models.Category  true  "Category data"
// @Success      201       {object}  models.Category  "Category created successfully"
// @Failure      400       {object}  map[string]interface{}  "Invalid input data"
// @Failure      500       {object}  map[string]interface{}  "Failed to create category"
// @Router       /categories [post]
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

// GetCategories godoc
// @Summary      Get Categories
// @Description  Retrieve a list of categories
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200       {object}  []models.Category  "Categories retrieved successfully"
// @Failure      500       {object}  map[string]interface{}  "Failed to retrieve categories"
// @Router       /categories [get]
func (cc *CategoryController) GetCategories(c *gin.Context) {
	categories, err := cc.Service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID godoc
// @Summary      Get Category by ID
// @Description  Retrieve a category by its ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Category ID"
// @Success      200   {object}  models.Category  "Category retrieved successfully"
// @Failure      404   {object}  map[string]interface{}  "Category not found"
// @Router       /categories/{id} [get]
func (cc *CategoryController) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := cc.Service.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory godoc
// @Summary      Update Category
// @Description  Update an existing category by its ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category  body      models.Category  true  "Updated category data"
// @Success      200       {object}  models.Category  "Category updated successfully"
// @Failure      400       {object}  map[string]interface{}  "Invalid input data"
// @Failure      500       {object}  map[string]interface{}  "Failed to update category"
// @Router       /categories/{id} [put]
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

// DeleteCategory godoc
// @Summary      Delete Category
// @Description  Delete a category by its ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      204  {object}  nil  "Category deleted successfully"
// @Failure      500  {object}  map[string]interface{}  "Failed to delete category"
// @Router       /categories/{id} [delete]
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := cc.Service.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
