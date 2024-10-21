package controllers

import (
	"net/http"

	"products-api/service"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	Service *service.CustomerService
}

// GetTopCustomers handles the request for fetching top customers based on total spent
func (rc *ReportController) GetTopCustomers(c *gin.Context) {
	customers, err := rc.Service.GetTopCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top customers"})
		return
	}
	c.JSON(http.StatusOK, customers)
}
