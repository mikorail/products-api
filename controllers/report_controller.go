package controllers

import (
	"net/http"

	"products-api/service"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	Service *service.CustomerService
}

// GetTopCustomers godoc
// @Summary      Get Top Customers
// @Description  Fetch a list of top customers based on their total spent amount
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Customer  "Top customers retrieved successfully"
// @Failure      500  {object}  map[string]interface{}  "Failed to fetch top customers"
// @Router       /reports/top-customers [get]
func (rc *ReportController) GetTopCustomers(c *gin.Context) {
	customers, err := rc.Service.GetTopCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top customers"})
		return
	}
	c.JSON(http.StatusOK, customers)
}
