package models

import "time"

type ProductBillerSummaryNotification struct {
	DateTime time.Time `json:"datetime"`
	Total    int       `json:"total"`
	Active   int       `json:"active"`
	Inactive int       `json:"inactive"`
}
