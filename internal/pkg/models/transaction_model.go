package models

// Transaction represents the structure of the transaction message.
type Transaction struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}
