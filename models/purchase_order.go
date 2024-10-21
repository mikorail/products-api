package models

type PurchaseOrderRequest struct {
	ProductId  uint `json:"product_id"`
	Quantity   int  `json:"quantity"`
	CustomerId uint `json:"customer_id"`
}
