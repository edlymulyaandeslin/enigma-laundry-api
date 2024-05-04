package entity

type TrxDetail struct {
	Id           int `json:"id"`
	TrxId        int `json:"trxId"`
	ProductId    int `json:"productId"`
	ProductPrice int `json:"productPrice"`
	Qty          int `json:"qty"`
}
