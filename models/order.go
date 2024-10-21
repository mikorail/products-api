package models

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	Product    Product   `json:"product" gorm:"foreignKey:ProductID"`
	CustomerID uint      `json:"customer_id"`
}

type OrderResponse struct {
	ID           uint      `json:"id"`
	CustomerName string    `json:"customer_name"` // Ensure this matches the alias used in the query
	ProductName  string    `json:"product_name"`
	Quantity     int       `json:"quantity"`
	TotalPrice   float64   `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"` // Adjust the type as needed
}
