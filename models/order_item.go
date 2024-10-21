package models

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Order     Order   `json:"order" gorm:"foreignKey:OrderID"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}
