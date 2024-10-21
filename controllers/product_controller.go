package controllers

import (
	"net/http"
	"products-api/helpers"
	"products-api/models"
	"products-api/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ProductController struct {
	Service *service.ProductService
}

// GetProductsWithTotalSold handles the request for fetching products with total sold quantities
func (pc *ProductController) GetProductsWithTotalSold(c *gin.Context) {
	// Read query parameters for pagination, filtering, sorting, and ordering
	page, _ := strconv.Atoi(c.Query("page")) // Current page number
	if page < 1 {
		page = 1 // Default to first page if not specified or invalid
	}

	limit, _ := strconv.Atoi(c.Query("limit")) // Number of items per page
	if limit < 1 {
		limit = 10 // Default limit
	}

	category := c.Query("category") // Optional filtering by category
	sortBy := c.Query("sort_by")    // Field to sort by (e.g., "total_sold_amount")
	order := c.Query("order")       // "asc" or "desc"

	// Call service to fetch products with total sold quantities, passing the parameters
	products, total, err := pc.Service.GetProductsWithTotalSold(page, limit, category, sortBy, order)
	if err != nil {
		helpers.RespondError(c, http.StatusInternalServerError, "Failed to fetch products", err)
		return
	}

	// Return the result along with total count for pagination
	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"products": products,
	})
}

// GetProductReport handles the request for fetching a product report for dashboards
func (pc *ProductController) GetProductReport(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize == 0 {
		pageSize = 10
	}

	filter := make(map[string]interface{})

	// Handle filtering
	if name := c.Query("name"); name != "" {
		filter["name"] = name
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		filter["category_id"] = categoryID
	}
	if priceMin := c.Query("price_min"); priceMin != "" {
		filter["price_min"] = priceMin
	}
	if priceMax := c.Query("price_max"); priceMax != "" {
		filter["price_max"] = priceMax
	}
	if sortField := c.Query("sort_field"); sortField != "" {
		filter["sort_field"] = sortField
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filter["sort_order"] = sortOrder
	}

	report, err := pc.Service.GetProductReport(page, pageSize, filter)
	if err != nil {
		helpers.RespondError(c, http.StatusInternalServerError, "Failed to fetch order history", err)
		return
	}

	helpers.RespondSuccess(c, "Order history retrieved successfully", report)
}

// CreateProduct handles the creation of a new product
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the product using the validator
	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.Service.CreateProduct(&product); err != nil {
		helpers.RespondError(c, http.StatusInternalServerError, "Failed to create product", err)
		return
	}

	helpers.RespondSuccess(c, "Product created successfully", product)
}

// GetProducts retrieves a paginated list of products
func (pc *ProductController) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize == 0 {
		pageSize = 10
	}

	products, err := pc.Service.GetProducts(page, pageSize)
	if err != nil {
		helpers.RespondError(c, http.StatusInternalServerError, "Failed to retrieve products", err)
		return
	}

	helpers.RespondSuccess(c, "Order history retrieved successfully", products)
}

// GetProductByID retrieves a product by ID
func (pc *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32) // Parse string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	strID := strconv.Itoa(int(idUint))
	product, err := pc.Service.GetProductByID(strID) // Convert to uint
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates a product by ID
// UpdateProduct updates a product by ID
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 32) // Parse string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product.ID = uint(idUint) // Assign converted value to product.ID

	// Validate the product using the validator
	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.Service.UpdateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by ID
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := pc.Service.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetProductsWithTotalSold handles the request for fetching products with total sold quantities
func (pc *ProductController) GetProductsWithTotalSoldReportCSV(c *gin.Context) {
	filter := make(map[string]interface{})

	// Handle filtering
	if name := c.Query("name"); name != "" {
		filter["name"] = name
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		filter["category_id"] = categoryID
	}
	if priceMin := c.Query("price_min"); priceMin != "" {
		filter["price_min"] = priceMin
	}
	if priceMax := c.Query("price_max"); priceMax != "" {
		filter["price_max"] = priceMax
	}
	if sortField := c.Query("sort_field"); sortField != "" {
		filter["sort_field"] = sortField
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filter["sort_order"] = sortOrder
	}
	// Call service to fetch products with total sold quantities, passing the parameters
	products, total, err := pc.Service.GetProductReportAsZip(filter)
	if err != nil {
		helpers.RespondError(c, http.StatusInternalServerError, "Failed to fetch products", err)
		return
	}

	// Return the result along with total count for pagination
	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"products": products,
	})
}

// GetProductsWithTotalSold handles the request for fetching products with total sold quantities
func (pc *ProductController) GetProductsByFilter(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize == 0 {
		pageSize = 10
	}

	filter := make(map[string]interface{})

	// Handle filtering
	if name := c.Query("name"); name != "" {
		filter["name"] = name
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		filter["category_id"] = categoryID
	}
	if priceMin := c.Query("price_min"); priceMin != "" {
		filter["price_min"] = priceMin
	}
	if priceMax := c.Query("price_max"); priceMax != "" {
		filter["price_max"] = priceMax
	}
	if sortField := c.Query("sort_field"); sortField != "" {
		filter["sort_field"] = sortField
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filter["sort_order"] = sortOrder
	}
	// Call service to fetch products with total sold quantities, passing the parameters
	products, err := pc.Service.GetProductsByFilter(filter, page, pageSize)
	if err != nil {
		helpers.RespondError(c, http.StatusInternalServerError, "Failed to fetch products", err)
		return
	}

	// Return the result along with total count for pagination
	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
