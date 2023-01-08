package models

type Transaction struct {
	Id              uint64 `json:"id"`
	TransactionName string `json:"transactionName"`
	Amount          string `json:"amount"`
	CreatedBy       string `json:"createdBy"`
}
