package repository

import (
	"fmt"
	"products-api/models"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	DB *gorm.DB
}

// GetTopCustomers retrieves the top 10 customers based on total amount spent.
func (cr *CustomerRepository) GetTopCustomers() ([]models.CustomerResponse, error) {
	var customers []models.CustomerResponse
	query := `
		SELECT c.id AS customer_id, c.name AS customer_name, SUM(o.total_price) AS total_spent
		FROM customers c
		JOIN orders o ON o.customer_id = c.id
		GROUP BY c.id
		ORDER BY total_spent DESC
		LIMIT 10
	`
	err := cr.DB.Raw(query).Scan(&customers).Error
	fmt.Println("customers :", customers)

	fmt.Println("AAAA")
	return customers, err
}

// Create adds a new customer to the database
func (pr *CustomerRepository) Create(customer *models.Customer) error {
	result := pr.DB.Create(customer)
	return result.Error
}

// GetAll retrieves all customers with pagination
func (cr *CustomerRepository) GetAll(page, pageSize int) ([]models.Customer, error) {
	var customers []models.Customer
	offset := (page - 1) * pageSize
	result := cr.DB.Offset(offset).Limit(pageSize).Find(&customers)
	return customers, result.Error
}

// GetByID retrieves a customer by ID
func (cr *CustomerRepository) GetByID(id string) (models.Customer, error) {
	var customer models.Customer
	result := cr.DB.First(&customer, id)
	if result.Error != nil {
		return customer, fmt.Errorf("customer not found")
	}
	return customer, nil
}

// Update modifies an existing customer
func (cr *CustomerRepository) Update(customer *models.Customer) error {
	result := cr.DB.Save(customer)
	return result.Error
}

// Delete removes a customer by ID
func (cr *CustomerRepository) Delete(id string) error {
	result := cr.DB.Delete(&models.Customer{}, id)
	return result.Error
}
