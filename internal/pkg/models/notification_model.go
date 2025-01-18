package models

import "time"

type CountProductBillerNotification struct {
	DateTime time.Time `json:"datetime"`
	Total    uint      `json:"total"`
	Active   uint      `json:"active"`
	Inactive uint      `json:"inactive"`
}
