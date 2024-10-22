// Package routes Products API.
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//	Host: localhost:8080
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta

package routes

import (
	"products-api/config"
	"products-api/controllers"
	"products-api/repository"
	"products-api/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	// Initialize repositories
	productRepo := &repository.ProductRepository{DB: config.DB}
	customerRepo := &repository.CustomerRepository{DB: config.DB}
	orderRepo := &repository.OrderRepository{DB: config.DB}
	categoryRepo := &repository.CategoryRepository{DB: config.DB}

	// Initialize services
	productService := &service.ProductService{ProductRepo: productRepo}
	customerService := &service.CustomerService{CustomerRepo: customerRepo}
	orderService := &service.OrderService{OrderRepo: orderRepo}
	categoryService := &service.CategoryService{CategoryRepo: categoryRepo}

	// Initialize controllers
	productController := &controllers.ProductController{Service: productService}
	reportController := &controllers.ReportController{Service: customerService}
	orderController := &controllers.OrderController{Service: orderService}
	categoryController := &controllers.CategoryController{Service: categoryService}
	customerController := &controllers.CustomerController{Service: customerService}

	// Products API with total sold quantities
	r.GET("/products/total-sold", productController.GetProductsWithTotalSold)

	// Export CSV
	r.GET("/products/total-sold-csv", productController.GetProductsWithTotalSoldReportCSV)

	// Top 10 customers by total spent
	r.GET("/customers/top", reportController.GetTopCustomers)

	// Order history with related products and customers
	r.GET("/orders/history", orderController.GetOrderHistory)

	// Product Report
	r.GET("/products/report", productController.GetProductReport)

	r.POST("/orders/purchase", orderController.PurchaseOrder)

	// Products routes
	r.POST("/products", productController.CreateProduct)
	r.GET("/products", productController.GetProducts)
	r.GET("/products/:id", productController.GetProductByID)
	r.PUT("/products/:id", productController.UpdateProduct)
	r.DELETE("/products/:id", productController.DeleteProduct)
	r.GET("/products/filter", productController.GetProductsByFilter)

	// Categories routes
	r.POST("/categories", categoryController.CreateCategory)
	r.GET("/categories", categoryController.GetCategories)
	r.GET("/categories/:id", categoryController.GetCategoryByID)
	r.PUT("/categories/:id", categoryController.UpdateCategory)
	r.DELETE("/categories/:id", categoryController.DeleteCategory)

	// Customers routes
	r.POST("/costumers", customerController.CreateCustomer)
	r.GET("/costumers", customerController.GetCustomers)
	r.GET("/costumers/:id", customerController.GetCustomerByID)
	r.PUT("/costumers/:id", customerController.UpdateCustomer)
	r.DELETE("/costumers/:id", customerController.DeleteCustomer)

	// Orders routes
	r.POST("/orders", orderController.CreateOrder)
	r.GET("/orders", orderController.GetOrders)
	// r.GET("/orders/:id", orderController.GetOrder)
	r.PUT("/orders/:id", orderController.UpdateOrder)
	r.DELETE("/orders/:id", orderController.DeleteOrder)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
