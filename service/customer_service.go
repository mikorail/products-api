package service

import (
	"fmt"
	"products-api/helpers"
	"products-api/models"
	"products-api/repository"
	"time"
)

type CustomerService struct {
	CustomerRepo *repository.CustomerRepository
}

// isCacheValid checks if the cached data is still valid based on some criteria.
func (cs *CustomerService) isCacheValid(customers []models.CustomerResponse, cacheTimestamp time.Time) bool {
	// Define cache expiry duration
	cacheExpiryDuration := time.Minute * 5 // Set the cache to expire after 5 minutes

	// Ensure there are customers in cache
	if len(customers) == 0 {
		return false // Cache is invalid if there are no customers
	}

	// Check if the cache has expired
	if time.Since(cacheTimestamp) > cacheExpiryDuration {
		return false // Cache is invalid due to expiry
	}

	// Further conditions can be added as per your requirements
	return true // Return true if all checks pass, indicating the cache is valid
}

// GetTopCustomers fetches the top customers and caches the result.
func (cs *CustomerService) GetTopCustomers() ([]models.CustomerResponse, error) {
	cacheKey := "top_customers_by_spent"
	var customers []models.CustomerResponse

	found, err := helpers.GetCache(cacheKey, &customers)
	if err != nil {
		fmt.Println("Error retrieving data from cache:", err)
	}

	// Check if the cached data is still valid
	cacheTimestamp := time.Now() // Use the current time for the cache check
	if found && cs.isCacheValid(customers, cacheTimestamp) && customers == nil {
		fmt.Println("Cache hit, returning cached data")
		return customers, nil
	}

	fmt.Println("Cache miss or invalid, fetching data from the database")
	fetchedCustomers, err := cs.CustomerRepo.GetTopCustomers()
	if err != nil {
		return nil, err
	}
	// Check for mismatches between the cache and the database
	if found && len(customers) != len(fetchedCustomers) {
		fmt.Println("Mismatch detected: cache product count does not match DB product count")
		// Optionally, you could decide to clear the cache here
		// helpers.DeleteCache(cacheKey)
	}

	// Set cache for 10 minutes after successfully fetching from DB
	err = helpers.SetCache(cacheKey, fetchedCustomers, 10*time.Minute)
	if err != nil {
		fmt.Println("Error setting cache:", err)
	}

	return fetchedCustomers, nil
}

// CreateCategory creates a new customer
func (cs *CustomerService) CreateCustomer(customer *models.Customer) error {
	return cs.CustomerRepo.Create(customer)
}

// GetProducts retrieves paginated customer
func (cs *CustomerService) GetCustomers(page, pageSize int) ([]models.Customer, error) {
	return cs.CustomerRepo.GetAll(page, pageSize)
}

// GetProductByID retrieves a customer by ID
func (cs *CustomerService) GetCostumerByID(id string) (models.Customer, error) {
	return cs.CustomerRepo.GetByID(id)
}

// UpdateCustomer updates a customer by ID
func (cs *CustomerService) UpdateCustomer(customer *models.Customer) error {
	return cs.CustomerRepo.Update(customer)
}

// DeleteCustomer deletes a customer by ID
func (cs *CustomerService) DeleteCustomer(id string) error {
	return cs.CustomerRepo.Delete(id)
}
