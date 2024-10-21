package models

type ProductReportResponse struct {
	TotalProducts int64     `json:"total_products"`
	TotalStock    int64     `json:"total_stock"`
	AveragePrice  float64   `json:"average_price"`
	Products      []Product `json:"products"`
}
