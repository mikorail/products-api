package repository

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"products-api/helpers"
	"products-api/models"
	"strconv"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func (pr *ProductRepository) GetProductsWithTotalSold(page, limit int, category, sortBy, order string) ([]models.ProductResponse, int, error) {
	var products []models.ProductResponse
	var total int64
	fmt.Println(page, limit, category, sortBy, order)
	// Base query for fetching products
	query := pr.DB.Table("products p"). // Use Debug() to print SQL
						Select("p.id AS product_id, p.name AS product_name, c.name AS category_name, COALESCE(CAST(SUM(oi.quantity) AS FLOAT), 0) AS total_sold_amount").
						Joins("JOIN categories c ON p.category_id = c.id").
						Joins("INNER JOIN order_items oi ON oi.product_id = p.id").
						Group("p.id, c.name")
	// Apply filtering
	if category != "" {
		query = query.Where("c.name = ?", category)
	}

	// Count total products for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sortBy != "" {
		if order == "desc" {
			query = query.Order(sortBy + " DESC")
		} else {
			query = query.Order(sortBy + " ASC")
		}
	}

	// Apply pagination
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Debug() will print the generated SQL query
	if err := query.Debug().Scan(&products).Error; err != nil {
		return nil, 0, err
	}
	fmt.Println("QUERY:", query.Statement.SQL)

	return products, int(total), nil
}

// GetProductReport retrieves a report of all products with filtering, sorting, and pagination.
func (pr *ProductRepository) GetProductReport(page, pageSize int, filter map[string]interface{}) (models.ProductReportResponse, error) {
	var report models.ProductReportResponse

	// Base query
	baseQuery := pr.DB.Model(&models.Product{}).
		Select("COUNT(*) AS total_products, SUM(stock_quantity) AS total_stock, AVG(price) AS average_price")

	// Apply filtering
	if name, ok := filter["name"]; ok {
		baseQuery = baseQuery.Where("name ILIKE ?", "%"+name.(string)+"%")
	}
	if categoryID, ok := filter["category_id"]; ok {
		baseQuery = baseQuery.Where("category_id = ?", categoryID)
	}
	if priceMin, ok := filter["price_min"]; ok {
		baseQuery = baseQuery.Where("price >= ?", priceMin)
	}
	if priceMax, ok := filter["price_max"]; ok {
		baseQuery = baseQuery.Where("price <= ?", priceMax)
	}

	// Execute the statistics query
	var totalProducts int64
	var totalStock int64
	var averagePrice float64

	// Use raw SQL for fetching the statistics
	statsQuery := `
		SELECT COUNT(*) AS total_products, SUM(stock_quantity) AS total_stock, AVG(price) AS average_price 
		FROM products
		WHERE 1=1`

	// Apply filtering conditions to the stats query
	if name, ok := filter["name"]; ok {
		statsQuery += " AND name ILIKE '%" + name.(string) + "%'"
	}
	if categoryID, ok := filter["category_id"]; ok {
		statsQuery += " AND category_id = " + fmt.Sprintf("%v", categoryID)
	}
	if priceMin, ok := filter["price_min"]; ok {
		statsQuery += " AND price >= " + fmt.Sprintf("%v", priceMin)
	}
	if priceMax, ok := filter["price_max"]; ok {
		statsQuery += " AND price <= " + fmt.Sprintf("%v", priceMax)
	}

	// Execute the raw stats query
	rows, err := pr.DB.Raw(statsQuery).Rows()
	if err != nil {
		return report, err
	}
	defer rows.Close() // Ensure rows are closed after scanning

	// Scan the results into variables
	if rows.Next() {
		if err := rows.Scan(&totalProducts, &totalStock, &averagePrice); err != nil {
			return report, err
		}
	}

	// Populate the report with statistics
	report.TotalProducts = totalProducts
	report.TotalStock = totalStock
	report.AveragePrice = averagePrice

	// Retrieve paginated product details
	productsQuery := pr.DB.Model(&models.Product{}).
		Preload("Category").
		Offset((page - 1) * pageSize).
		Limit(pageSize)

	// Apply sorting
	if sortField, ok := filter["sort_field"]; ok {
		sortOrder := "ASC"
		if order, ok := filter["sort_order"]; ok {
			sortOrder = order.(string)
		}
		productsQuery = productsQuery.Order(sortField.(string) + " " + sortOrder)
	}

	// Execute the products query and scan the results
	var products []models.Product
	if err := productsQuery.Find(&products).Error; err != nil {
		return report, err
	}

	// Assign the products to the report
	report.Products = products

	return report, nil
}

// GetProductReport retrieves a report of all products with filtering, sorting, and pagination.
func (pr *ProductRepository) GetProducts(filter map[string]interface{}) (models.ProductReportResponse, error) {
	var report models.ProductReportResponse
	var totalProducts int64
	var totalStock int64
	var averagePrice float64

	// Use raw SQL for fetching the statistics
	statsQuery := `
		SELECT COUNT(*) AS total_products, SUM(stock_quantity) AS total_stock, CAST(AVG(price) AS NUMERIC(10, 2)) AS average_price
		FROM products
		WHERE 1=1`

	// Apply filtering conditions to the stats query
	if name, ok := filter["name"]; ok {
		statsQuery += " AND name ILIKE '%" + name.(string) + "%'"
	}
	if categoryID, ok := filter["category_id"]; ok {
		statsQuery += " AND category_id = " + fmt.Sprintf("%v", categoryID)
	}
	if priceMin, ok := filter["price_min"]; ok {
		statsQuery += " AND price >= " + fmt.Sprintf("%v", priceMin)
	}
	if priceMax, ok := filter["price_max"]; ok {
		statsQuery += " AND price <= " + fmt.Sprintf("%v", priceMax)
	}

	// Execute the raw stats query

	fmt.Println("QUERY:", statsQuery)

	rows, err := pr.DB.Raw(statsQuery).Rows()
	if err != nil {
		return report, err
	}

	defer rows.Close() // Ensure rows are closed after scanning

	// Scan the results into variables
	if rows.Next() {
		if err := rows.Scan(&totalProducts, &totalStock, &averagePrice); err != nil {
			return report, err
		}
	}

	// Populate the report with statistics
	report.TotalProducts = totalProducts
	report.TotalStock = totalStock
	// report.AveragePrice = averagePrice

	// Creating the products

	var products []models.Product

	// Use raw SQL for fetching the statistics
	productsQuery := `
		select  
			p.id as product_id,
			p."name" as product_name,
			p.description as product_description,
			p.price as product_price, 
			c."name" as category_name,
			p.stock_quantity as stock,
			p.is_active as is_active, 
			p.created_at 
		 from products p join categories c on p.category_id = c.id 
		 where 1=1 `

	// Apply filtering conditions to the stats query
	if name, ok := filter["name"]; ok {
		productsQuery += " AND name ILIKE '%" + name.(string) + "%'"
	}
	if categoryID, ok := filter["category_id"]; ok {
		productsQuery += " AND category_id = " + fmt.Sprintf("%v", categoryID)
	}
	if priceMin, ok := filter["price_min"]; ok {
		productsQuery += " AND price >= " + fmt.Sprintf("%v", priceMin)
	}
	if priceMax, ok := filter["price_max"]; ok {
		productsQuery += " AND price <= " + fmt.Sprintf("%v", priceMax)
	}

	// Execute the raw stats query

	fmt.Println("QUERY:", productsQuery)

	productsRows, err := pr.DB.Raw(productsQuery).Rows()
	if err != nil {
		return report, err
	}
	defer productsRows.Close() // Ensure rows are closed after scanning

	// Loop through each row and scan into the products slice and increment price
	var totalPrice float64

	for productsRows.Next() {
		var product models.Product
		var categoryName string // Assuming product model has a Category field or you need to assign categoryName manually

		// Scan the row into the product and categoryName variables
		err := productsRows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&categoryName,
			&product.StockQuantity,
			&product.IsActive,
			&product.CreatedAt,
		)
		if err != nil {
			return report, err
		}
		// Increment the price
		totalPrice += product.Price

		// Assuming the product model has a Category field (adjust as needed)
		product.Category = models.Category{
			Name: categoryName,
		}

		// Append the product to the slice
		products = append(products, product)
	}
	// Average calculation
	report.AveragePrice = totalPrice / float64(report.TotalProducts)

	// Assign the products to the report
	report.Products = products

	return report, nil
}

// Create adds a new product to the database
func (pr *ProductRepository) Create(product *models.Product) error {
	result := pr.DB.Create(product)
	return result.Error
}

// GetAll retrieves all products with pagination
func (pr *ProductRepository) GetAll(page, pageSize int) ([]models.Product, error) {
	var products []models.Product
	offset := (page - 1) * pageSize
	result := pr.DB.Offset(offset).Limit(pageSize).Find(&products)
	return products, result.Error
}

// GetByID retrieves a product by ID
func (pr *ProductRepository) GetByID(id string) (models.Product, error) {
	var product models.Product
	result := pr.DB.First(&product, id)
	if result.Error != nil {
		return product, fmt.Errorf("product not found")
	}
	return product, nil
}

// Update modifies an existing product
func (pr *ProductRepository) Update(product *models.Product) error {
	result := pr.DB.Save(product)
	return result.Error
}

// Delete removes a product by ID
func (pr *ProductRepository) Delete(id string) error {
	result := pr.DB.Delete(&models.Product{}, id)
	return result.Error
}

// Helper function to ensure the "exports" directory exists
func ensureExportsDirExists() error {
	exportsDir := "exports"
	if _, err := os.Stat(exportsDir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := os.Mkdir(exportsDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create exports directory: %w", err)
		}
	}
	return nil
}

// Function to create and compress the CSV file without concurrency
func (pr *ProductRepository) GetProductReportAsZip(filter map[string]interface{}) ([]byte, string, error) {
	report, err := pr.GetProducts(filter)
	if err != nil {
		return nil, "", err
	}

	fmt.Println(filter)
	// Generate the filename based on the filter and timestamp
	fileName := helpers.GenerateFileName(filter)
	csvFileName := fileName + ".csv"
	zipFileName := fileName + ".zip"

	// Ensure the "exports" directory exists
	if err := ensureExportsDirExists(); err != nil {
		return nil, "", err
	}

	// Create the full path for the CSV file in the "exports" folder
	csvFilePath := filepath.Join("exports", csvFileName)

	// Open the CSV file for writing
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer csvFile.Close()

	// Create a CSV writer
	csvWriter := csv.NewWriter(csvFile)

	// Sequentially write the filter section to the CSV file
	csvWriter.Write([]string{"Filter Name", "Filter Value"})

	// Write the filter values
	for key, value := range filter {
		csvWriter.Write([]string{key, fmt.Sprintf("%v", value)})
	}

	// Separator row
	csvWriter.Write([]string{""})

	averagePriceStr := strconv.FormatFloat(report.AveragePrice, 'f', 2, 64) // 2 decimal places
	fmt.Println("averagePriceStr", averagePriceStr)
	// Write the totals section to the CSV file
	csvWriter.Write([]string{"Total Products", "Total Stock", "Average Price"})
	csvWriter.Write([]string{
		strconv.FormatInt(report.TotalProducts, 10),
		strconv.FormatInt(report.TotalStock, 10),
		fmt.Sprintf("%v", report.AveragePrice),
	})

	// Separator row
	csvWriter.Write([]string{""})

	// Write the product headers to the CSV file
	csvWriter.Write([]string{"Product ID", "Name", "Category", "Price", "Stock Quantity"})

	// Write product rows in sequence
	for _, product := range report.Products {
		strProductID := strconv.FormatUint(uint64(product.ID), 10)
		strStockQuantity := strconv.FormatUint(uint64(product.StockQuantity), 10)
		productRow := []string{
			strProductID,
			product.Name,
			product.Category.Name, // Assuming Product has a Category relationship
			fmt.Sprintf("%.2f", product.Price),
			strStockQuantity,
		}
		csvWriter.Write(productRow)
	}

	// Flush the CSV writer
	csvWriter.Flush()

	// Close the CSV file
	if err := csvFile.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close CSV file: %w", err)
	}

	// Now create the zip file in the same "exports" folder
	zipFilePath := filepath.Join("exports", zipFileName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(zipFile)

	// Add the CSV file to the zip archive
	csvInZip, err := zipWriter.Create(csvFileName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create entry for CSV in zip file: %w", err)
	}

	// Read the CSV file and write it into the zip archive
	csvData, err := os.ReadFile(csvFilePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read CSV file: %w", err)
	}

	if _, err := csvInZip.Write(csvData); err != nil {
		return nil, "", fmt.Errorf("failed to write CSV data to zip: %w", err)
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close zip writer: %w", err)
	}

	// Delete the CSV file after zip is successfully created
	if err := os.Remove(csvFilePath); err != nil {
		return nil, "", fmt.Errorf("failed to delete CSV file: %w", err)
	}

	// Return the zip file path and name
	return []byte(zipFilePath), zipFileName, nil
}

// GetProductReport retrieves a report of all products with filtering, sorting, and pagination.
func (pr *ProductRepository) GetProductsByFilter(filter map[string]interface{}, page, pageSize int) ([]models.Product, error) {
	var products []models.Product

	// Use raw SQL for fetching the statistics
	productsQuery := `
		select  
			p.id as product_id,
			p."name" as product_name,
			p.description as product_description,
			p.price as product_price, 
			c."name" as category_name,
			p.stock_quantity as stock,
			p.is_active as is_active, 
			p.created_at 
		 from products p join categories c on p.category_id = c.id 
		 where 1=1 `

	// Apply filtering conditions to the stats query
	if name, ok := filter["name"]; ok {
		productsQuery += " AND name ILIKE '%" + name.(string) + "%'"
	}
	if categoryID, ok := filter["category_id"]; ok {
		productsQuery += " AND category_id = " + fmt.Sprintf("%v", categoryID)
	}
	if priceMin, ok := filter["price_min"]; ok {
		productsQuery += " AND price >= " + fmt.Sprintf("%v", priceMin)
	}
	if priceMax, ok := filter["price_max"]; ok {
		productsQuery += " AND price <= " + fmt.Sprintf("%v", priceMax)
	}

	// Implement pagination using LIMIT and OFFSET
	offset := (page - 1) * pageSize
	productsQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	// Execute the raw stats query

	fmt.Println("QUERY:", productsQuery)

	productsRows, err := pr.DB.Raw(productsQuery).Rows()
	if err != nil {
		return products, err
	}
	defer productsRows.Close() // Ensure rows are closed after scanning

	// Loop through each row and scan into the products slice and increment price
	var totalPrice float64

	for productsRows.Next() {
		var product models.Product
		var categoryName string // Assuming product model has a Category field or you need to assign categoryName manually

		// Scan the row into the product and categoryName variables
		err := productsRows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&categoryName,
			&product.StockQuantity,
			&product.IsActive,
			&product.CreatedAt,
		)
		if err != nil {
			return products, err
		}
		// Increment the price
		totalPrice += product.Price

		// Assuming the product model has a Category field (adjust as needed)
		product.Category = models.Category{
			Name: categoryName,
		}

		// Append the product to the slice
		products = append(products, product)
	}
	// Average calculation

	return products, nil
}
