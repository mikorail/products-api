package controllers

import (
	"fmt"
	"net/http"
	"products-api/helpers"
	"products-api/models"
	"products-api/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	Service *service.OrderService
}

// GetOrderHistory handles the request for fetching order history
func (oc *OrderController) GetOrderHistory(c *gin.Context) {
	orders, err := oc.Service.GetOrderHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order history"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// CreateOrder handles the creation of a new order
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := oc.Service.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// PurchaseOrderRequest handles the creation of a new order
func (oc *OrderController) PurchaseOrder(c *gin.Context) {
	fmt.Println("params : ", c.Request)

	var order models.PurchaseOrderRequest
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := oc.Service.PurchaseOrderRequest(&order); err != nil {
		helpers.RespondError(c, http.StatusInternalServerError, "Failed to create orders", err)
		return
	}

	helpers.RespondSuccess(c, "Orders created successfully successfully", order)
}

// GetProducts retrieves a paginated list of costumers
func (oc *OrderController) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize == 0 {
		pageSize = 10
	}

	customers, err := oc.Service.GetOrders(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// GetOrderByID retrieves a order by ID
func (oc *OrderController) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	order, err := oc.Service.GetCostumerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder updates a order by ID
func (oc *OrderController) UpdateOrder(c *gin.Context) {
	var costumer models.Order
	if err := c.ShouldBindJSON(&costumer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := oc.Service.UpdateOrder(&costumer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, costumer)
}

// DeleteOrder deletes a costumer by ID
func (oc *OrderController) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := oc.Service.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
