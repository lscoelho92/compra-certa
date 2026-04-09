package dto

import "compra-certa/api/internal/models"

type CreateProductRequest struct {
	Name         string   `json:"name"`
	Description  *string  `json:"description"`
	CategoryID   *string  `json:"categoryId"`
	DefaultPrice *float64 `json:"defaultPrice"`
}

type UpdateProductRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	CategoryID   *string  `json:"categoryId"`
	DefaultPrice *float64 `json:"defaultPrice"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type ProductDetailResponse struct {
	Product      models.Product        `json:"product"`
	PriceHistory []models.ProductPrice `json:"priceHistory"`
}

type CreatePurchaseRequest struct {
	PurchaseDate *string               `json:"purchaseDate"`
	Items        []PurchaseItemRequest `json:"items"`
}

type UpdatePurchaseRequest struct {
	PurchaseDate *string                `json:"purchaseDate"`
	Items        *[]PurchaseItemRequest `json:"items"`
}

type PurchaseItemRequest struct {
	ProductID string  `json:"productId"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
}

type PurchaseItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Category    string  `json:"category"`
}

type PurchaseResponse struct {
	ID           string                 `json:"id"`
	PurchaseDate string                 `json:"purchase_date"`
	Month        string                 `json:"month"`
	TotalPrice   float64                `json:"total_price"`
	Items        []PurchaseItemResponse `json:"items"`
}

type PurchaseListItem = PurchaseResponse
