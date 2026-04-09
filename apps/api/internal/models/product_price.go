package models

import "time"

type ProductPrice struct {
    ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    ProductID string    `json:"productId" gorm:"type:uuid"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `json:"createdAt"`
}

func (ProductPrice) TableName() string {
    return "product_prices"
}
