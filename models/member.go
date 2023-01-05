package models

type Member struct {
	PhoneNumber    uint64 `json:"phone_number"`
	AmountToBePaid uint64 `json:"amount_to_be_paid"`
	TransactionId  uint64 `json:"transaction_id"`
}
