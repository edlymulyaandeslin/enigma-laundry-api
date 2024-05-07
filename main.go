package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"submission-project-enigma-laundry/config"
	"submission-project-enigma-laundry/entity"

	"github.com/gin-gonic/gin"
)

var db = config.ConnectDb()

func main() {
	// Tulis kode kamu disini
	defer db.Close()

	router := gin.Default()

	// api customers
	customerRouter := router.Group("/customers")
	{
		customerRouter.GET("/", getAllCustomer)
		customerRouter.GET("/:id", getCustomerById)
		customerRouter.POST("/", createCustomer)
		customerRouter.PUT("/:id", updatedCustomer)
		customerRouter.DELETE("/:id", deletedCustomer)
	}

	// api employee
	employeeRouter := router.Group("/employees")
	{
		employeeRouter.GET("/", getAllEmployee)
		employeeRouter.GET("/:id", getEmployeeById)
		employeeRouter.POST("/", createEmployee)
		employeeRouter.PUT("/:id", updatedEmployee)
		employeeRouter.DELETE("/:id", deletedEmployee)
	}

	// api products
	productRouter := router.Group("/products")
	{
		productRouter.GET("/", getAllProduct)
		productRouter.GET("/:id", getProductById)
		productRouter.POST("/", createProduct)
		productRouter.PUT("/:id", updatedProduct)
		productRouter.DELETE("/:id", deletedProduct)
	}

	// api transaction
	transactionRouter := router.Group("/transactions")
	{
		transactionRouter.GET("/", getAllTrx)
		transactionRouter.GET("/:id_bill", getTrxById)
		transactionRouter.POST("/", transaction)
	}

	router.Run(":3000")
}

// form yang harus di isi user
type Form struct {
	BillDate   string `json:"billDate"`
	EntryDate  string `json:"entryDate"`
	FinishDate string `json:"finishDate"`
	EmployeeId int    `json:"employeeId"`
	CustomerId int    `json:"customerId"`
	ProductId  int    `json:"productId"`
	Qty        int    `json:"qty"`
}

// transaction function
// get all transaction
func getAllTrx(c *gin.Context) {

	entryDate := c.Query("entry_date")
	finishDate := c.Query("finish_date")
	productName := c.Query("product_name")

	fmt.Println("product name", productName)

	query := "SELECT transactions.id, bill_date, entry_date, finish_date  , employe_id, customer_id FROM transactions"

	var rows *sql.Rows
	var err error

	if entryDate != "" {
		query += " WHERE entry_date ILIKE '%' || $1 || '%'"
		rows, err = db.Query(query, entryDate)
	} else if finishDate != "" {
		query += " WHERE finish_date ILIKE '%' || $1 || '%'"
		rows, err = db.Query(query, finishDate)
	} else if productName != "" {
		query += " JOIN trx_detail ON transactions.id = trx_detail.trx_id JOIN mst_product ON trx_detail.product_id = mst_product.id WHERE mst_product.name = $1"
		rows, err = db.Query(query, productName)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	var transactions = []gin.H{}
	for rows.Next() {
		var transaction entity.Transaction
		err := rows.Scan(&transaction.Id, &transaction.BillDate, &transaction.EntryDate, &transaction.FinishDate, &transaction.EmployeeId, &transaction.CustomerId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		var employee entity.Employee
		query2 := "SELECT * FROM mst_employee WHERE id = $1"
		db.QueryRow(query2, transaction.EmployeeId).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address)

		var customer entity.Customer
		query3 := "SELECT * FROM mst_customer WHERE id = $1"
		db.QueryRow(query3, transaction.CustomerId).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)

		query4 := "SELECT * FROM trx_detail WHERE trx_id = $1"
		rows2, err := db.Query(query4, transaction.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer rows2.Close()

		var trxDetail []entity.TrxDetail
		total := 0
		for rows2.Next() {
			var detail entity.TrxDetail
			err := rows2.Scan(&detail.Id, &detail.TrxId, &detail.ProductId, &detail.ProductPrice, &detail.Qty)
			total += detail.ProductPrice * detail.Qty
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			trxDetail = append(trxDetail, detail)
		}

		// transactions = append(transactions, transaction)
		data := gin.H{
			"id":          transaction.Id,
			"billDate":    transaction.BillDate,
			"entryDate":   transaction.EntryDate,
			"finishDate":  transaction.FinishDate,
			"employee":    employee,
			"customer":    customer,
			"billDetails": trxDetail,
			"totalBill":   total,
		}

		transactions = append(transactions, data)
	}
	if len(transactions) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Get all transactions success",
			"data":    transactions,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transactions not found"})
	}
}

// get transaction by id
func getTrxById(c *gin.Context) {
	id := c.Param("id_bill")

	billId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	// find transaction
	var transaction entity.Transaction
	query := "SELECT * FROM transactions WHERE id = $1"
	err = db.QueryRow(query, billId).Scan(&transaction.Id, &transaction.BillDate, &transaction.EntryDate, &transaction.FinishDate, &transaction.EmployeeId, &transaction.CustomerId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Transactions not found",
		})
		return
	}

	var employee entity.Employee
	query2 := "SELECT * FROM mst_employee WHERE id = $1"
	db.QueryRow(query2, transaction.EmployeeId).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address)

	var customer entity.Customer
	query3 := "SELECT * FROM mst_customer WHERE id = $1"
	db.QueryRow(query3, transaction.CustomerId).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)

	query4 := "SELECT * FROM trx_detail WHERE trx_id = $1"
	rows, err := db.Query(query4, billId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	var trxDetail []entity.TrxDetail
	total := 0
	for rows.Next() {
		var detail entity.TrxDetail
		err := rows.Scan(&detail.Id, &detail.TrxId, &detail.ProductId, &detail.ProductPrice, &detail.Qty)
		total += detail.ProductPrice * detail.Qty
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		trxDetail = append(trxDetail, detail)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get transaction by id success",
		"data": gin.H{
			"id":          transaction.Id,
			"billDate":    transaction.BillDate,
			"entryDate":   transaction.EntryDate,
			"finishDate":  transaction.FinishDate,
			"employee":    employee,
			"customer":    customer,
			"billDetails": trxDetail,
			"totalBill":   total,
		},
	})
}

// new transaction
func transaction(c *gin.Context) {
	var newTransaction entity.Transaction
	var newTrxDetail entity.TrxDetail

	var form Form

	tx, err := db.Begin()

	if err = c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTransaction.BillDate = form.BillDate
	newTransaction.EntryDate = form.EntryDate
	newTransaction.FinishDate = form.FinishDate
	newTransaction.EmployeeId = form.EmployeeId
	newTransaction.CustomerId = form.CustomerId

	newTrxDetail.ProductId = form.ProductId
	newTrxDetail.Qty = form.Qty

	// create newTrx
	trxId := createTrx(tx, newTransaction, c)

	// create detail newTrx
	trxDetailId, productId := createDetailTrx(tx, newTrxDetail, trxId, c)

	// update Transaction
	productPrice := updateTrx(tx, trxDetailId, productId, c)

	err = tx.Commit()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "telor"})
		return
	} else {
		newTransaction.Id = trxId
		newTrxDetail.Id = trxDetailId

		c.JSON(http.StatusCreated, gin.H{
			"message": "New transaction created",
			"data": gin.H{
				"id":         newTransaction.Id,
				"billDate":   newTransaction.BillDate,
				"entryDate":  newTransaction.EntryDate,
				"finishDate": newTransaction.FinishDate,
				"employeeId": newTransaction.EmployeeId,
				"customerId": newTransaction.CustomerId,
				"billDetail": []gin.H{
					{
						"id":           newTrxDetail.Id,
						"billId":       trxId,
						"productId":    newTrxDetail.ProductId,
						"productPrice": productPrice,
						"qty":          newTrxDetail.Qty,
					},
				},
			},
		})
	}
}
func createTrx(tx *sql.Tx, newTransaction entity.Transaction, c *gin.Context) int {
	query1 := "INSERT INTO transactions (bill_date, entry_date, finish_date, employe_id, customer_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var trxId int
	err := tx.QueryRow(query1, newTransaction.BillDate, newTransaction.EntryDate, newTransaction.FinishDate, newTransaction.EmployeeId, newTransaction.CustomerId).Scan(&trxId)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	return trxId
}
func createDetailTrx(tx *sql.Tx, newTrxDetail entity.TrxDetail, trxId int, c *gin.Context) (int, int) {
	// insert transaction detail
	query2 := "INSERT INTO trx_detail (trx_id, product_id, product_price, qty) VALUES ($1, $2, $3, $4) RETURNING id, product_id"

	var trxDetailId, product_id int
	productPrice := 0
	err := tx.QueryRow(query2, trxId, newTrxDetail.ProductId, productPrice, newTrxDetail.Qty).Scan(&trxDetailId, &product_id)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	return trxDetailId, product_id
}
func updateTrx(tx *sql.Tx, trxDetailId int, productId int, c *gin.Context) int {
	// find product
	queryFindProduct := "SELECT * FROM mst_product WHERE id = $1"

	var product entity.Product
	err := tx.QueryRow(queryFindProduct, productId).Scan(&product.Id, &product.Name, &product.Price, &product.Unit)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}

	// update price detail trx
	productPrice := product.Price
	query3 := "UPDATE trx_detail SET product_price = $1 WHERE id = $2"
	_, err = tx.Exec(query3, productPrice, trxDetailId)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	return productPrice
}

// customer function
func getAllCustomer(c *gin.Context) {
	query := "SELECT * FROM mst_customer"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	customers := []entity.Customer{}
	for rows.Next() {
		var customer entity.Customer
		err := rows.Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		customers = append(customers, customer)
	}

	if len(customers) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "get all customer success",
			"data":    customers,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
	}
}

func getCustomerById(c *gin.Context) {
	id := c.Param("id")

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	query := "SELECT * FROM mst_customer WHERE id =  $1"

	customer := entity.Customer{}

	err = db.QueryRow(query, customerId).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Customer not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get customer by id success",
		"data":    customer,
	})
}

func createCustomer(c *gin.Context) {
	newCustomer := entity.Customer{}

	err := c.ShouldBind(&newCustomer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryInsert := "INSERT INTO mst_customer (name, phone_number, address) VALUES ($1, $2, $3) RETURNING id"

	var customerId int
	err = db.QueryRow(queryInsert, newCustomer.Name, newCustomer.PhoneNumber, newCustomer.Address).Scan(&customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	newCustomer.Id = customerId
	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer created successfully",
		"data":    newCustomer,
	})
}

func updatedCustomer(c *gin.Context) {
	id := c.Param("id")

	customerId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	findCSQuery := "SELECT * FROM mst_customer WHERE id = $1"
	var customer entity.Customer
	err = db.QueryRow(findCSQuery, customerId).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)

	// cek apakah customer ditemukan atau tidak
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Customer not found",
		})
		return
	}

	var updateCustomer entity.Customer
	err = c.ShouldBind(&updateCustomer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(updateCustomer.Name) == "" {
		updateCustomer.Name = customer.Name
	}
	if strings.TrimSpace(updateCustomer.PhoneNumber) == "" {
		updateCustomer.PhoneNumber = customer.PhoneNumber
	}
	if strings.TrimSpace(updateCustomer.Address) == "" {
		updateCustomer.Address = customer.Address
	}

	queryUpdate := "UPDATE mst_customer SET name = $1, phone_number = $2, address = $3 WHERE id = $4"
	_, err = db.Exec(queryUpdate, updateCustomer.Name, updateCustomer.PhoneNumber, updateCustomer.Address, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update customer",
		})
	} else {
		updateCustomer.Id = customerId
		c.JSON(http.StatusOK, gin.H{
			"message": "Customer updated successfully",
			"data":    updateCustomer,
		})
	}
}

func deletedCustomer(c *gin.Context) {
	id := c.Param("id")

	customerId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	findCSQuery := "SELECT * FROM mst_customer WHERE id = $1"
	var customer entity.Customer
	err = db.QueryRow(findCSQuery, customerId).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)

	// cek apakah customer ditemukan atau tidak
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Customer not found",
		})
		return
	}

	queryDelete := "DELETE FROM mst_customer WHERE id = $1"
	_, err = db.Exec(queryDelete, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete customer",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Customer deleted successfully",
		})
	}
}

// employee function
func getAllEmployee(c *gin.Context) {
	query := "SELECT * FROM mst_employee"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	employees := []entity.Employee{}
	for rows.Next() {
		var employee entity.Employee
		err := rows.Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		employees = append(employees, employee)
	}

	if len(employees) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "get all employee success",
			"data":    employees,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Employee not found"})
	}
}

func getEmployeeById(c *gin.Context) {
	id := c.Param("id")

	employeeId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	query := "SELECT * FROM mst_employee WHERE id =  $1"

	employee := entity.Employee{}

	err = db.QueryRow(query, employeeId).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Employee not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get employee by id success",
		"data":    employee,
	})
}

func createEmployee(c *gin.Context) {
	newEmployee := entity.Employee{}

	err := c.ShouldBind(&newEmployee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryInsert := "INSERT INTO mst_employee (name, phone_number, address) VALUES ($1, $2, $3) RETURNING id"

	var employeeId int
	err = db.QueryRow(queryInsert, newEmployee.Name, newEmployee.PhoneNumber, newEmployee.Address).Scan(&employeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	newEmployee.Id = employeeId
	c.JSON(http.StatusCreated, gin.H{
		"message": "Employee created successfully",
		"data":    newEmployee,
	})
}

func updatedEmployee(c *gin.Context) {
	id := c.Param("id")

	EmployeeId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	findEMQuery := "SELECT * FROM mst_employee WHERE id = $1"
	var employee entity.Employee
	err = db.QueryRow(findEMQuery, EmployeeId).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address)

	// cek apakah employee ditemukan atau tidak
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Employee not found",
		})
		return
	}

	var employeeUpdate entity.Employee
	err = c.ShouldBind(&employeeUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(employeeUpdate.Name) == "" {
		employeeUpdate.Name = employee.Name
	}
	if strings.TrimSpace(employeeUpdate.PhoneNumber) == "" {
		employeeUpdate.PhoneNumber = employee.PhoneNumber
	}
	if strings.TrimSpace(employeeUpdate.Address) == "" {
		employeeUpdate.Address = employee.Address
	}

	queryUpdate := "UPDATE mst_employee SET name = $1, phone_number = $2, address = $3 WHERE id = $4"
	_, err = db.Exec(queryUpdate, employeeUpdate.Name, employeeUpdate.PhoneNumber, employeeUpdate.Address, EmployeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update employee",
		})
	} else {
		employeeUpdate.Id = EmployeeId
		c.JSON(http.StatusOK, gin.H{
			"message": "Employee updated successfully",
			"data":    employeeUpdate,
		})
	}
}

func deletedEmployee(c *gin.Context) {
	id := c.Param("id")

	employeeId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	findEMQuery := "SELECT * FROM mst_employee WHERE id = $1"
	var employee entity.Employee
	err = db.QueryRow(findEMQuery, employeeId).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address)

	// cek apakah employee ditemukan atau tidak
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Employee not found",
		})
		return
	}

	queryDelete := "DELETE FROM mst_employee WHERE id = $1"
	_, err = db.Exec(queryDelete, employeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete employee",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Employee deleted successfully",
		})
	}
}

// product function
func getAllProduct(c *gin.Context) {
	productName := c.Query("product_name")

	query := "SELECT * FROM mst_product"

	var rows *sql.Rows
	var err error

	if productName != "" {
		query += " WHERE name ILIKE '%' || $1 || '%'"
		rows, err = db.Query(query, productName)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	products := []entity.Product{}
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		products = append(products, product)
	}

	if len(products) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "get all products success",
			"data":    products,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
	}
}

func getProductById(c *gin.Context) {
	id := c.Param("id")

	productId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	query := "SELECT * FROM mst_product WHERE id =  $1"

	product := entity.Product{}
	err = db.QueryRow(query, productId).Scan(&product.Id, &product.Name, &product.Price, &product.Unit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get product by id success",
		"data":    product,
	})
}

func createProduct(c *gin.Context) {
	newProduct := entity.Product{}

	err := c.ShouldBind(&newProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryInsert := "INSERT INTO mst_product (name, price, unit) VALUES ($1, $2, $3) RETURNING id"

	var productId int
	err = db.QueryRow(queryInsert, newProduct.Name, newProduct.Price, newProduct.Unit).Scan(&productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	newProduct.Id = productId
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    newProduct,
	})
}

func updatedProduct(c *gin.Context) {
	id := c.Param("id")

	productId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	findProductQuery := "SELECT * FROM mst_product WHERE id = $1"
	var product entity.Product
	err = db.QueryRow(findProductQuery, productId).Scan(&product.Id, &product.Name, &product.Price, &product.Unit)

	// cek apakah product ditemukan atau tidak
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	var productUpdate entity.Product
	err = c.ShouldBind(&productUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(productUpdate.Name) == "" {
		productUpdate.Name = product.Name
	}
	if productUpdate.Price == 0 {
		productUpdate.Price = product.Price
	}
	if strings.TrimSpace(productUpdate.Unit) == "" {
		productUpdate.Unit = product.Unit
	}

	queryUpdate := "UPDATE mst_product SET name = $1, price = $2, unit = $3 WHERE id = $4"
	_, err = db.Exec(queryUpdate, productUpdate.Name, productUpdate.Price, productUpdate.Unit, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product",
		})
	} else {
		productUpdate.Id = productId
		c.JSON(http.StatusOK, gin.H{
			"message": "Product updated successfully",
			"data":    productUpdate,
		})
	}
}

func deletedProduct(c *gin.Context) {
	id := c.Param("id")

	productId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	findProductQuery := "SELECT * FROM mst_product WHERE id = $1"
	var product entity.Product
	err = db.QueryRow(findProductQuery, productId).Scan(&product.Id, &product.Name, &product.Price, &product.Unit)

	// cek apakah product ditemukan atau tidak
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	queryDelete := "DELETE FROM mst_product WHERE id = $1"
	_, err = db.Exec(queryDelete, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete product",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Product deleted successfully",
		})
	}
}
