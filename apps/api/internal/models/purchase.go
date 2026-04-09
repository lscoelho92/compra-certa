package models

import "time"

type Purchase struct {
	ID          string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PurchasedAt time.Time `json:"purchasedAt"`
	Total       float64   `json:"total"`
}

func (Purchase) TableName() string {
	return "purchases"
}
