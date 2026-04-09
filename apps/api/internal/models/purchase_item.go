package models

type PurchaseItem struct {
	ID         string  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PurchaseID string  `json:"purchaseId" gorm:"type:uuid"`
	ProductID  string  `json:"productId" gorm:"type:uuid"`
	Quantity   float64 `json:"quantity"`
	Price      float64 `json:"price"`
}

func (PurchaseItem) TableName() string {
	return "purchase_items"
}
