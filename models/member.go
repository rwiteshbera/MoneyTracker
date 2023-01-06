package models

type Member struct {
	PhoneNumber    uint64 `json:"phone_number"`
	TransactionId  uint64 `json:"transaction_id"`
	AmountToBePaid uint64 `json:"amount_to_be_paid"`
}
