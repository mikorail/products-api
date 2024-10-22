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

// GetOrderHistory godoc
// @Summary      Get Order History
// @Description  Fetch the order history for the user
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Order  "Order history retrieved successfully"
// @Failure      500  {object}  map[string]interface{}  "Failed to fetch order history"
// @Router       /orders/history [get]
func (oc *OrderController) GetOrderHistory(c *gin.Context) {
	orders, err := oc.Service.GetOrderHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order history"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// CreateOrder godoc
// @Summary      Create Order
// @Description  Create a new order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order  body      models.Order  true  "Order data"
// @Success      201    {object}  models.Order  "Order created successfully"
// @Failure      400    {object}  map[string]interface{}  "Invalid order data"
// @Failure      500    {object}  map[string]interface{}  "Failed to create order"
// @Router       /orders [post]
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

// PurchaseOrder godoc
// @Summary      Purchase Order
// @Description  Create a new purchase order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        purchaseOrder  body      models.PurchaseOrderRequest  true  "Purchase order data"
// @Success      200            {object}  models.PurchaseOrderRequest  "Purchase order created successfully"
// @Failure      400            {object}  map[string]interface{}  "Invalid purchase order data"
// @Failure      500            {object}  map[string]interface{}  "Failed to create purchase order"
// @Router       /orders/purchase [post]
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

// GetOrders godoc
// @Summary      Get Orders
// @Description  Retrieve a paginated list of orders
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        pageSize   query     int     false  "Number of items per page (default: 10)"
// @Success      200        {object}  map[string]interface{}  "Orders retrieved successfully"
// @Failure      500        {object}  map[string]interface{}  "Failed to retrieve orders"
// @Router       /orders [get]
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

// GetOrderByID godoc
// @Summary      Get Order by ID
// @Description  Retrieve an order by its ID
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  models.Order  "Order retrieved successfully"
// @Failure      404  {object}  map[string]interface{}  "Order not found"
// @Router       /orders/{id} [get]
func (oc *OrderController) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	order, err := oc.Service.GetCostumerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder godoc
// @Summary      Update Order
// @Description  Update an existing order by its ID
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order  body      models.Order  true  "Updated order data"
// @Success      200    {object}  models.Order  "Order updated successfully"
// @Failure      400    {object}  map[string]interface{}  "Invalid input data"
// @Failure      500    {object}  map[string]interface{}  "Failed to update order"
// @Router       /orders/{id} [put]
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

// DeleteOrder godoc
// @Summary      Delete Order
// @Description  Delete an order by its ID
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      204  {object}  nil  "Order deleted successfully"
// @Failure      500  {object}  map[string]interface{}  "Failed to delete order"
// @Router       /orders/{id} [delete]
func (oc *OrderController) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := oc.Service.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
