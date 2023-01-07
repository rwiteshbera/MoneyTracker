package models

type Member struct {
	PhoneNumber    string `json:"phone_number"`
	TransactionId  uint64 `json:"transaction_id"`
	AmountToBePaid string `json:"amount_to_be_paid"`
	CreatedBy      string `json:"created_by"`
}
