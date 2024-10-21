package models

import "time"

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomerResponse struct {
	CustomerId   int     `json:"customer_id" gorm:"customer_id"`
	CustomerName string  `json:"customer_name" gorm:"customer_name"`
	TotalSpent   float64 `json:"total_spent"`
}
