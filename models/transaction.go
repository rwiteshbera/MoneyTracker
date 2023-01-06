package models

type Transaction struct {
	Id        uint64 `json:"id"`
	Amount    string `json:"amount"`
	CreatedBy string `json:"createdBy"`
}
