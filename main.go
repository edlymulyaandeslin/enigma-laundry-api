package main

import (
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

	customerRouter := router.Group("/customers")
	{
		customerRouter.GET("/", getAllCustomer)
		customerRouter.GET("/:id", getCustomerById)
		customerRouter.POST("/", createCustomer)
		customerRouter.PUT("/:id", updatedCustomer)
		customerRouter.DELETE("/:id", deletedCustomer)
	}

	employeeRouter := router.Group("/employees")
	{
		employeeRouter.GET("/", getAllEmployee)
		employeeRouter.GET("/:id", getEmployeeById)
		employeeRouter.POST("/", createEmployee)
		employeeRouter.PUT("/:id", updatedEmployee)
		employeeRouter.DELETE("/:id", deletedEmployee)
	}

	router.Run(":3000")
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

	findCSQuery := "SELECT * FROM mst_employee WHERE id = $1"
	var employee entity.Employee
	err = db.QueryRow(findCSQuery, employeeId).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address)

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
