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

// CreateCustomer godoc
// @Summary      Create Customer
// @Description  Create a new customer
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        customer  body      models.Customer  true  "Customer data"
// @Success      201       {object}  models.Customer  "Customer created successfully"
// @Failure      400       {object}  map[string]interface{}  "Invalid input data"
// @Failure      500       {object}  map[string]interface{}  "Failed to create customer"
// @Router       /customers [post]
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

// GetCustomers godoc
// @Summary      Get Customers
// @Description  Retrieve a paginated list of customers
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        pageSize   query     int     false  "Number of items per page (default: 10)"
// @Success      200        {object}  map[string]interface{}  "Customers retrieved successfully"
// @Failure      500        {object}  map[string]interface{}  "Failed to retrieve customers"
// @Router       /customers [get]
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

// GetCustomerByID godoc
// @Summary      Get Customer by ID
// @Description  Retrieve a customer by its ID
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Customer ID"
// @Success      200  {object}  models.Customer  "Customer retrieved successfully"
// @Failure      404  {object}  map[string]interface{}  "Customer not found"
// @Router       /customers/{id} [get]
func (cc *CustomerController) GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	customer, err := cc.Service.GetCostumerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// UpdateCustomer godoc
// @Summary      Update Customer
// @Description  Update an existing customer by its ID
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        customer  body      models.Customer  true  "Updated customer data"
// @Success      200       {object}  models.Customer  "Customer updated successfully"
// @Failure      400       {object}  map[string]interface{}  "Invalid input data"
// @Failure      500       {object}  map[string]interface{}  "Failed to update customer"
// @Router       /customers/{id} [put]
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

// DeleteCustomer godoc
// @Summary      Delete Customer
// @Description  Delete a customer by its ID
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Customer ID"
// @Success      204  {object}  nil  "Customer deleted successfully"
// @Failure      500  {object}  map[string]interface{}  "Failed to delete customer"
// @Router       /customers/{id} [delete]
func (cc *CustomerController) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := cc.Service.DeleteCustomer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
