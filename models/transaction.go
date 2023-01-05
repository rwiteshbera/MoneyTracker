package models

type Transaction struct {
	Id        uint64 `json:"id"`
	Amount    uint64 `json:"amount"`
	CreatedBy uint64 `json:"createdBy"`
}
