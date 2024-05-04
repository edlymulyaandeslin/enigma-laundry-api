package entity

type Transaction struct {
	Id         int    `json:"id"`
	BillDate   string `json:"billDate"`
	EntryDate  string `json:"entryDate"`
	FinishDate string `json:"finishDate"`
	EmployeeId int    `json:"employeId"`
	CustomerId int    `json:"customerId"`
}
