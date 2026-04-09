package models

import "time"

type Product struct {
    ID           string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Name         string    `json:"name"`
    Description  *string   `json:"description"`
    CategoryID   *string   `json:"categoryId" gorm:"type:uuid"`
    DefaultPrice *float64  `json:"defaultPrice"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}
