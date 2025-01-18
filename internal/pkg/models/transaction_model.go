package models

// Transaction represents the structure of the transaction message.
type Transaction struct {
	ID        uint   `json:"id"`
	PartnerID uint   `json:"partner_id"`
	ProductID uint   `json:"product_id"`
	BillerID  uint   `json:"biller_id"`
	Status    string `json:"status"`
}
