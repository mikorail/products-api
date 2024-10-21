package repository

import (
	"fmt"
	"products-api/models"
	"strconv"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

// PurchaseOrder handles the logic for purchasing an order
func (or *OrderRepository) PurchaseOrder(por *models.PurchaseOrderRequest) error {
	fmt.Println("Inside PurchaseOrder")

	// Begin a transaction
	tx := or.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	// Ensure that the transaction is rolled back in case of an error
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
			if r != nil {
				fmt.Println("Recovered in PurchaseOrder:", r)
			}
		}
	}()

	// Convert ProductId to string
	productIDStr := strconv.FormatUint(uint64(por.ProductId), 10)

	// Raw SQL query to retrieve product by ID
	var prodCheck models.Product
	query := `
		SELECT id, name, price, stock_quantity, is_active 
		FROM products 
		WHERE id = ?`

	// Execute the raw query
	if err := tx.Raw(query, productIDStr).Scan(&prodCheck).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("product is not found: %v", err)
	}
	fmt.Println("Product check passed")

	// Check if the product is active
	if !prodCheck.IsActive {
		tx.Rollback()
		return fmt.Errorf("product is not active")
	} else if prodCheck.StockQuantity < por.Quantity {
		tx.Rollback()
		return fmt.Errorf("requested quantity exceeds available stock")
	}

	// Check customer existence
	var customerCheck models.Customer
	costumerIDStr := strconv.FormatUint(uint64(por.CustomerId), 10)
	queryCustomer := `
		SELECT id
		FROM customers 
		WHERE id = ?`

	// Execute the raw query
	if err := tx.Raw(queryCustomer, costumerIDStr).Scan(&customerCheck).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("customer is not found: %v", err)
	}
	fmt.Println("Customer check passed")

	// Create a new order
	newOrder := models.Order{
		ProductID:  por.ProductId,
		CustomerID: por.CustomerId,
		Quantity:   por.Quantity,
		TotalPrice: float64(por.Quantity) * prodCheck.Price,
	}

	OrdID, err := or.CreateShowID(&newOrder) // Pass transaction to CreateShowID
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("can't create order: %v", err)
	}

	newOrderItems := models.OrderItem{
		ProductID: por.ProductId,
		OrderID:   OrdID,
		Quantity:  por.Quantity,
		UnitPrice: prodCheck.Price,
	}

	if err := or.CreateOrderItems(&newOrderItems); err != nil { // Pass transaction to CreateOrderItems
		tx.Rollback()
		return fmt.Errorf("can't create order items: %v", err)
	}

	fmt.Println("Order and order items created successfully")

	calcStock := prodCheck.StockQuantity - por.Quantity

	// Update product stock
	queryUpdate := `
		UPDATE products 
		SET stock_quantity = ? 
		WHERE id = ?`

	// Execute the raw query
	if err := tx.Exec(queryUpdate, calcStock, productIDStr).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update product stock: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// GetOrderHistory retrieves order history with related products and orders.
func (or *OrderRepository) GetOrderHistory() ([]models.OrderResponse, error) {
	var orders []models.OrderResponse

	query := `
		SELECT
			o.id, 
			c.name AS customer_name,
			p.name AS product_name, 
			oi.quantity, 
			o.total_price, 
			o.created_at
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
		JOIN order_items oi ON oi.order_id = o.id
		JOIN products p ON oi.product_id = p.id
	`

	// Execute the raw query and scan results into the orders slice
	err := or.DB.Raw(query).Scan(&orders).Error

	if err != nil {
		return nil, err // Handle error and return nil if any error occurs
	}

	return orders, nil
}

// Create adds a new order to the database
func (pr *OrderRepository) Create(order *models.Order) error {
	result := pr.DB.Create(order)
	return result.Error
}

func (pr *OrderRepository) CreateShowID(order *models.Order) (id uint, err error) {
	result := pr.DB.Create(order)
	return order.ID, result.Error
}

// GetAll retrieves all orders with pagination
func (or *OrderRepository) GetAll(page, pageSize int) ([]models.Order, error) {
	var orders []models.Order
	offset := (page - 1) * pageSize
	result := or.DB.Offset(offset).Limit(pageSize).Find(&orders)
	return orders, result.Error
}

// GetByID retrieves a order by ID
func (or *OrderRepository) GetByID(id string) (models.Order, error) {
	var order models.Order
	result := or.DB.First(&order, id)
	if result.Error != nil {
		return order, fmt.Errorf("order not found")
	}
	return order, nil
}

// Update modifies an existing order
func (or *OrderRepository) Update(order *models.Order) error {
	result := or.DB.Save(order)
	return result.Error
}

// Delete removes a order by ID
func (or *OrderRepository) Delete(id string) error {
	result := or.DB.Delete(&models.Order{}, id)
	return result.Error
}

func (pr *OrderRepository) CreateOrderItems(orderItems *models.OrderItem) error {
	result := pr.DB.Create(orderItems)
	return result.Error
}
