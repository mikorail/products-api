package service

import (
	"time"

	"products-api/helpers"
	"products-api/models"
	"products-api/repository"
)

type OrderService struct {
	OrderRepo    *repository.OrderRepository
	ProductRepo  *repository.ProductRepository
	CustomerRepo *repository.CustomerRepository
}

// PurchaseOrderRequest makes and order item and orders a order by ID
func (os *OrderService) PurchaseOrderRequest(por *models.PurchaseOrderRequest) error {
	return os.OrderRepo.PurchaseOrder(por)
}

// GetOrderHistory fetches order history and caches the result.
func (os *OrderService) GetOrderHistory() ([]models.OrderResponse, error) {
	cacheKey := "order_history"
	var orders []models.OrderResponse

	if found, _ := helpers.GetCache(cacheKey, &orders); found {
		return orders, nil
	}

	orders, err := os.OrderRepo.GetOrderHistory()
	if err != nil {
		return nil, err
	}

	// Set cache for 10 minutes
	helpers.SetCache(cacheKey, orders, 10*time.Minute)

	return orders, nil
}

// CreateCategory creates a new order
func (os *OrderService) CreateOrder(order *models.Order) error {
	return os.OrderRepo.Create(order)
}

// GetProducts retrieves paginated order
func (os *OrderService) GetOrders(page, pageSize int) ([]models.Order, error) {
	return os.OrderRepo.GetAll(page, pageSize)
}

// GetProductByID retrieves a order by ID
func (os *OrderService) GetCostumerByID(id string) (models.Order, error) {
	return os.OrderRepo.GetByID(id)
}

// UpdateOrder updates a order by ID
func (os *OrderService) UpdateOrder(order *models.Order) error {
	return os.OrderRepo.Update(order)
}

// DeleteOrder deletes a order by ID
func (os *OrderService) DeleteOrder(id string) error {
	return os.OrderRepo.Delete(id)
}
