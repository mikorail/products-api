package service

import (
	"fmt"
	"time"

	"products-api/helpers"
	"products-api/models"
	"products-api/repository"
)

type ProductService struct {
	ProductRepo *repository.ProductRepository
}

// GetProductReport fetches the product report with filtering, sorting, and pagination.
func (ps *ProductService) GetProductReport(page, pageSize int, filter map[string]interface{}) (models.ProductReportResponse, error) {
	return ps.ProductRepo.GetProductReport(page, pageSize, filter)
}

// GetProductsWithTotalSold fetches products with total sold quantities and caches the result.
func (ps *ProductService) GetProductsWithTotalSold(page, limit int, category, sortBy, order string) ([]models.ProductResponse, int, error) {
	// Generate a unique cache key based on parameters to avoid conflicts
	cacheKey := fmt.Sprintf("products_with_total_sold:%d:%d:%s:%s:%s", page, limit, category, sortBy, order)

	var products []models.ProductResponse
	var total int

	// Try to fetch data from cache
	found, err := helpers.GetCache(cacheKey, &products)
	if err != nil {
		fmt.Println("Error retrieving data from cache:", err)
	}

	// Check if the cached data is still valid
	cacheTimestamp := time.Now() // Use the current time for the cache check
	if found && ps.isCacheValid(products, page, limit, category, sortBy, order, cacheTimestamp) {
		fmt.Println("Cache hit, returning cached data")
		return products, total, nil
	}

	// If cache not found or invalid, fetch data from the database
	fmt.Println("Cache miss or invalid, fetching data from the database")
	fetchedProducts, total, err := ps.ProductRepo.GetProductsWithTotalSold(page, limit, category, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	// Check for mismatches between the cache and the database
	if found && len(products) != len(fetchedProducts) {
		fmt.Println("Mismatch detected: cache product count does not match DB product count")
		// Optionally, you could decide to clear the cache here
		// helpers.DeleteCache(cacheKey)
	}

	// Set cache for 10 minutes after successfully fetching from DB
	err = helpers.SetCache(cacheKey, fetchedProducts, 10*time.Minute)
	if err != nil {
		fmt.Println("Error setting cache:", err)
	}

	return fetchedProducts, total, nil
}

// GetProductsWithTotalSold fetches products with total sold quantities and caches the result.
func (ps *ProductService) GetProductReportAsZip(filter map[string]interface{}) ([]byte, string, error) {

	zip, name, error := ps.ProductRepo.GetProductReportAsZip(filter)
	if error != nil {
		return nil, "error", error
	}

	return zip, name, nil
}

// CreateProduct creates a new product
func (ps *ProductService) CreateProduct(product *models.Product) error {
	fmt.Println("product: ", product)
	return ps.ProductRepo.Create(product)
}

// GetProducts retrieves paginated products
func (ps *ProductService) GetProducts(page, pageSize int) ([]models.Product, error) {
	return ps.ProductRepo.GetAll(page, pageSize)
}

func (ps *ProductService) GetProductsByFilter(filter map[string]interface{}, page, pageSize int) ([]models.Product, error) {
	return ps.ProductRepo.GetProductsByFilter(filter, page, pageSize)
}

// GetProductByID retrieves a product by ID
func (ps *ProductService) GetProductByID(id string) (models.Product, error) {
	return ps.ProductRepo.GetByID(id)
}

// UpdateProduct updates a product by ID
func (ps *ProductService) UpdateProduct(product *models.Product) error {
	return ps.ProductRepo.Update(product)
}

// DeleteProduct deletes a product by ID
func (ps *ProductService) DeleteProduct(id string) error {
	return ps.ProductRepo.Delete(id)
}

// isCacheValid checks if the cached data is still valid based on some criteria.
func (ps *ProductService) isCacheValid(products []models.ProductResponse, page, limit int, category, sortBy, order string, cacheTimestamp time.Time) bool {
	// Define cache expiry duration
	cacheExpiryDuration := time.Minute * 5 // Set the cache to expire after 5 minutes

	// Check if the cache has expired
	if time.Since(cacheTimestamp) > cacheExpiryDuration {
		return false // Cache is invalid due to expiry
	}

	// Check if category filter is applied and no products found
	if category != "" && len(products) == 0 {
		return false // Cache is invalid if there's a category filter and no products found
	}

	// Additional validation: Compare product count with expected count based on pagination
	expectedCount := limit
	if page > 1 {
		expectedCount = (page - 1) * limit // Adjust based on pagination
	}

	if len(products) < expectedCount {
		return false // Cache is invalid if the number of cached products is less than expected
	}

	// Further conditions can be added as per your requirements
	return true // Return true if all checks pass, indicating the cache is valid
}
