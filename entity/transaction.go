package entity

type Transaction struct {
	Id         int    `json:"id"`
	BillDate   string `json:"billDate"`
	EntryDate  string `json:"entryDate"`
	FinishDate string `json:"finishDate"`
	EmployeeId int    `json:"employeeId"`
	CustomerId int    `json:"customerId"`
}
