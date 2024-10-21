package models

import (
	"time"
)

type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `json:"name" binding:"required"`
	Description   string    `json:"description"`
	Price         float64   `json:"price" binding:"required,gt=0"`
	CategoryID    uint      `json:"category_id" binding:"required,gt=0"`
	StockQuantity int       `json:"stock_quantity" binding:"required,gt=0"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Category      Category  `gorm:"foreignKey:CategoryID" json:"category"`
}

type ProductResponse struct {
	ProductID       uint    `json:"product_id" gorm:"product_id"`
	ProductName     string  `json:"product_name" gorm:"product_name"`
	CategoryName    string  `json:"category_name" gorm:"category_name"`
	TotalSoldAmount float64 `json:"total_sold_amount" gorm:"total_sold_amount"`
}
