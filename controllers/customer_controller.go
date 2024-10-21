package controllers

import (
	"net/http"
	"strconv"

	"products-api/models"
	"products-api/service"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	Service *service.CustomerService
}

// CreateCustomer handles the creation of a new customer
func (cc *CustomerController) CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Service.CreateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetProducts retrieves a paginated list of costumers
func (cc *CustomerController) GetCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize == 0 {
		pageSize = 10
	}

	customers, err := cc.Service.GetCustomers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// GetCustomerByID retrieves a customer by ID
func (cc *CustomerController) GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	customer, err := cc.Service.GetCostumerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// UpdateCustomer updates a customer by ID
func (cc *CustomerController) UpdateCustomer(c *gin.Context) {
	var costumer models.Customer
	if err := c.ShouldBindJSON(&costumer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Service.UpdateCustomer(&costumer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, costumer)
}

// DeleteCustomer deletes a costumer by ID
func (cc *CustomerController) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := cc.Service.DeleteCustomer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
