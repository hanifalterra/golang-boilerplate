package models

// Transaction represents the structure of the transaction message.
type Transaction struct {
	ID        int    `json:"id"`
	PartnerID int    `json:"partner_id"`
	ProductID int    `json:"product_id"`
	BillerID  int    `json:"biller_id"`
	Status    string `json:"status"`
}
