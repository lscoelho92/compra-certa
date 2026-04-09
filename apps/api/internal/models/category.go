package models

import "time"

type Category struct {
    ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
}

func (Category) TableName() string {
    return "categories"
}
