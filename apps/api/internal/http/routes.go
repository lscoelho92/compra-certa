package httpapi

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"compra-certa/api/internal/http/dto"
	"compra-certa/api/internal/models"
)

// BuildVersion is injected at build time via -ldflags. Defaults to v1.
var BuildVersion = "v1"

func getAPIVersion() string {
	version := strings.TrimSpace(BuildVersion)
	if version == "" {
		return "v1"
	}
	return version
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		MaxAge:           0,
	}))

	apiVersion := getAPIVersion()
	router.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "API Version: "+apiVersion)
	})

	loadPurchaseResponse := func(purchaseID string) (dto.PurchaseResponse, error) {
		type purchaseRow struct {
			PurchaseID  string
			PurchasedAt time.Time
			Total       float64
			ProductID   string
			ProductName string
			Quantity    float64
			UnitPrice   float64
			Category    *string
		}

		var rows []purchaseRow
		if err := db.Table("purchases").
			Select("purchases.id as purchase_id, purchases.purchased_at as purchased_at, purchases.total as total, purchase_items.product_id as product_id, products.name as product_name, purchase_items.quantity as quantity, purchase_items.price as unit_price, categories.name as category").
			Joins("join purchase_items on purchase_items.purchase_id = purchases.id").
			Joins("join products on products.id = purchase_items.product_id").
			Joins("left join categories on categories.id = products.category_id").
			Where("purchases.id = ?", purchaseID).
			Order("purchase_items.id asc").
			Scan(&rows).Error; err != nil {
			return dto.PurchaseResponse{}, err
		}

		if len(rows) == 0 {
			return dto.PurchaseResponse{}, gorm.ErrRecordNotFound
		}

		date := rows[0].PurchasedAt.Format("2006-01-02")
		month := ""
		if len(date) >= 7 {
			month = date[:7]
		}

		items := make([]dto.PurchaseItemResponse, 0, len(rows))
		for _, row := range rows {
			category := ""
			if row.Category != nil {
				category = *row.Category
			}
			items = append(items, dto.PurchaseItemResponse{
				ProductID:   row.ProductID,
				ProductName: row.ProductName,
				Quantity:    row.Quantity,
				UnitPrice:   row.UnitPrice,
				Category:    category,
			})
		}

		return dto.PurchaseResponse{
			ID:           rows[0].PurchaseID,
			PurchaseDate: date,
			Month:        month,
			TotalPrice:   rows[0].Total,
			Items:        items,
		}, nil
	}

	registerRoutes := func(r gin.IRoutes) {
		r.GET("/categories", func(c *gin.Context) {
			var categories []models.Category
			if err := db.Order("name asc").Find(&categories).Error; err != nil {
				respondError(c, err)
				return
			}

			c.JSON(http.StatusOK, categories)
		})

		r.POST("/categories", func(c *gin.Context) {
			var request dto.CreateCategoryRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				respondBadRequest(c, "invalid request body")
				return
			}

			request.Name = strings.TrimSpace(request.Name)
			if request.Name == "" {
				respondBadRequest(c, "name is required")
				return
			}

			category := models.Category{Name: request.Name}
			if err := db.Omit("id").Create(&category).Error; err != nil {
				respondError(c, err)
				return
			}

			c.JSON(http.StatusCreated, category)
		})

		r.DELETE("/categories/:id", func(c *gin.Context) {
			categoryID := c.Param("id")
			if strings.TrimSpace(categoryID) == "" {
				respondBadRequest(c, "missing category id")
				return
			}

			var count int64
			if err := db.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
				respondError(c, err)
				return
			}

			if count > 0 {
				respondBadRequest(c, "category has products")
				return
			}

			result := db.Delete(&models.Category{}, "id = ?", categoryID)
			if result.Error != nil {
				respondError(c, result.Error)
				return
			}

			if result.RowsAffected == 0 {
				respondNotFound(c, "category not found")
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		r.GET("/products", func(c *gin.Context) {
			var products []models.Product
			if err := db.Find(&products).Error; err != nil {
				respondError(c, err)
				return
			}

			c.JSON(http.StatusOK, products)
		})

		r.GET("/products/:id", func(c *gin.Context) {
			productID := c.Param("id")
			if strings.TrimSpace(productID) == "" {
				respondBadRequest(c, "missing product id")
				return
			}

			var product models.Product
			if err := db.First(&product, "id = ?", productID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					respondNotFound(c, "product not found")
					return
				}
				respondError(c, err)
				return
			}

			var history []models.ProductPrice
			if err := db.Order("created_at desc").Find(&history, "product_id = ?", productID).Error; err != nil {
				respondError(c, err)
				return
			}

			c.JSON(http.StatusOK, dto.ProductDetailResponse{
				Product:      product,
				PriceHistory: history,
			})
		})

		r.POST("/products", func(c *gin.Context) {
			var request dto.CreateProductRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				respondBadRequest(c, "invalid request body")
				return
			}

			request.Name = strings.TrimSpace(request.Name)
			if request.Name == "" {
				respondBadRequest(c, "name is required")
				return
			}

			var product models.Product
			if err := db.Transaction(func(tx *gorm.DB) error {
				product = models.Product{
					Name:         request.Name,
					Description:  request.Description,
					CategoryID:   request.CategoryID,
					DefaultPrice: request.DefaultPrice,
				}

				if err := tx.Omit("id").Create(&product).Error; err != nil {
					return err
				}

				if request.DefaultPrice != nil {
					price := models.ProductPrice{
						ProductID: product.ID,
						Price:     *request.DefaultPrice,
					}
					if err := tx.Omit("id").Create(&price).Error; err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				respondError(c, err)
				return
			}

			c.JSON(http.StatusCreated, product)
		})

		r.PATCH("/products/:id", func(c *gin.Context) {
			productID := c.Param("id")
			if strings.TrimSpace(productID) == "" {
				respondBadRequest(c, "missing product id")
				return
			}

			var request dto.UpdateProductRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				respondBadRequest(c, "invalid request body")
				return
			}

			if request.Name == nil && request.Description == nil && request.CategoryID == nil && request.DefaultPrice == nil {
				respondBadRequest(c, "no fields to update")
				return
			}

			if request.Name != nil {
				trimmed := strings.TrimSpace(*request.Name)
				if trimmed == "" {
					respondBadRequest(c, "name cannot be empty")
					return
				}
				request.Name = &trimmed
			}

			var product models.Product
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.First(&product, "id = ?", productID).Error; err != nil {
					return err
				}

				updates := map[string]interface{}{}
				if request.Name != nil {
					updates["name"] = *request.Name
				}
				if request.Description != nil {
					updates["description"] = *request.Description
				}
				if request.CategoryID != nil {
					updates["category_id"] = request.CategoryID
				}
				if request.DefaultPrice != nil {
					updates["default_price"] = request.DefaultPrice
				}

				if len(updates) == 0 {
					return nil
				}

				if err := tx.Model(&product).Updates(updates).Error; err != nil {
					return err
				}

				priceChanged := false
				if request.DefaultPrice != nil {
					if product.DefaultPrice == nil || *product.DefaultPrice != *request.DefaultPrice {
						priceChanged = true
					}
				}

				if priceChanged {
					price := models.ProductPrice{
						ProductID: product.ID,
						Price:     *request.DefaultPrice,
					}
					if err := tx.Omit("id").Create(&price).Error; err != nil {
						return err
					}
				}

				if request.Name != nil {
					product.Name = *request.Name
				}
				if request.Description != nil {
					product.Description = request.Description
				}
				if request.CategoryID != nil {
					product.CategoryID = request.CategoryID
				}
				if request.DefaultPrice != nil {
					product.DefaultPrice = request.DefaultPrice
				}
				return nil
			}); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					respondNotFound(c, "product not found")
					return
				}
				respondError(c, err)
				return
			}

			c.JSON(http.StatusOK, product)
		})

		r.DELETE("/products/:id", func(c *gin.Context) {
			productID := c.Param("id")
			if strings.TrimSpace(productID) == "" {
				respondBadRequest(c, "missing product id")
				return
			}

			result := db.Delete(&models.Product{}, "id = ?", productID)
			if result.Error != nil {
				if strings.Contains(result.Error.Error(), "violates foreign key constraint") {
					respondBadRequest(c, "product is referenced")
					return
				}
				respondError(c, result.Error)
				return
			}

			if result.RowsAffected == 0 {
				respondNotFound(c, "product not found")
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		r.GET("/purchases", func(c *gin.Context) {
			type purchaseRow struct {
				PurchaseID  string
				PurchasedAt time.Time
				Total       float64
				ProductID   string
				ProductName string
				Quantity    float64
				UnitPrice   float64
				Category    *string
			}

			var rows []purchaseRow
			if err := db.Table("purchases").
				Select("purchases.id as purchase_id, purchases.purchased_at as purchased_at, purchases.total as total, purchase_items.product_id as product_id, products.name as product_name, purchase_items.quantity as quantity, purchase_items.price as unit_price, categories.name as category").
				Joins("join purchase_items on purchase_items.purchase_id = purchases.id").
				Joins("join products on products.id = purchase_items.product_id").
				Joins("left join categories on categories.id = products.category_id").
				Order("purchases.purchased_at desc").
				Scan(&rows).Error; err != nil {
				respondError(c, err)
				return
			}

			responses := make([]dto.PurchaseResponse, 0)
			index := make(map[string]int)
			for _, row := range rows {
				idx, exists := index[row.PurchaseID]
				if !exists {
					date := row.PurchasedAt.Format("2006-01-02")
					month := ""
					if len(date) >= 7 {
						month = date[:7]
					}
					responses = append(responses, dto.PurchaseResponse{
						ID:           row.PurchaseID,
						PurchaseDate: date,
						Month:        month,
						TotalPrice:   row.Total,
						Items:        []dto.PurchaseItemResponse{},
					})
					idx = len(responses) - 1
					index[row.PurchaseID] = idx
				}

				category := ""
				if row.Category != nil {
					category = *row.Category
				}
				responses[idx].Items = append(responses[idx].Items, dto.PurchaseItemResponse{
					ProductID:   row.ProductID,
					ProductName: row.ProductName,
					Quantity:    row.Quantity,
					UnitPrice:   row.UnitPrice,
					Category:    category,
				})
			}

			c.JSON(http.StatusOK, responses)
		})

		r.POST("/purchases", func(c *gin.Context) {
			var request dto.CreatePurchaseRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				respondBadRequest(c, "invalid request body")
				return
			}

			if len(request.Items) == 0 {
				respondBadRequest(c, "items is required")
				return
			}

			productIDs := make([]string, 0, len(request.Items))
			seen := make(map[string]struct{})
			for index, item := range request.Items {
				item.ProductID = strings.TrimSpace(item.ProductID)
				if item.ProductID == "" {
					respondBadRequest(c, "productId is required")
					return
				}
				if item.Quantity <= 0 {
					respondBadRequest(c, "quantity must be greater than zero")
					return
				}
				if item.UnitPrice < 0 {
					respondBadRequest(c, "unitPrice must be zero or greater")
					return
				}
				request.Items[index] = item
				if _, ok := seen[item.ProductID]; !ok {
					seen[item.ProductID] = struct{}{}
					productIDs = append(productIDs, item.ProductID)
				}
			}

			var purchasedAt time.Time
			if request.PurchaseDate != nil && strings.TrimSpace(*request.PurchaseDate) != "" {
				parsed, err := time.ParseInLocation("2006-01-02", *request.PurchaseDate, time.Local)
				if err != nil {
					respondBadRequest(c, "purchaseDate must be in YYYY-MM-DD format")
					return
				}
				purchasedAt = parsed
			} else {
				purchasedAt = time.Now()
			}

			var count int64
			if err := db.Model(&models.Product{}).Where("id in ?", productIDs).Count(&count).Error; err != nil {
				respondError(c, err)
				return
			}
			if count != int64(len(productIDs)) {
				respondNotFound(c, "product not found")
				return
			}

			var purchase models.Purchase
			total := 0.0
			for _, item := range request.Items {
				total += item.Quantity * item.UnitPrice
			}
			if err := db.Transaction(func(tx *gorm.DB) error {
				purchase = models.Purchase{
					PurchasedAt: purchasedAt,
					Total:       total,
				}

				if err := tx.Omit("id").Create(&purchase).Error; err != nil {
					return err
				}

				for _, item := range request.Items {
					entry := models.PurchaseItem{
						PurchaseID: purchase.ID,
						ProductID:  item.ProductID,
						Quantity:   item.Quantity,
						Price:      item.UnitPrice,
					}
					if err := tx.Omit("id").Create(&entry).Error; err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				respondError(c, err)
				return
			}

			response, err := loadPurchaseResponse(purchase.ID)
			if err != nil {
				respondError(c, err)
				return
			}

			c.JSON(http.StatusCreated, response)
		})

		r.PATCH("/purchases/:id", func(c *gin.Context) {
			purchaseID := strings.TrimSpace(c.Param("id"))
			if purchaseID == "" {
				respondBadRequest(c, "missing purchase id")
				return
			}

			errProductNotFound := errors.New("product not found")

			var request dto.UpdatePurchaseRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				respondBadRequest(c, "invalid request body")
				return
			}

			if request.PurchaseDate == nil && request.Items == nil {
				respondBadRequest(c, "no fields to update")
				return
			}
			if request.PurchaseDate != nil && strings.TrimSpace(*request.PurchaseDate) == "" {
				respondBadRequest(c, "purchaseDate cannot be empty")
				return
			}

			if request.Items != nil && len(*request.Items) == 0 {
				respondBadRequest(c, "items must have at least one entry")
				return
			}

			productIDs := make([]string, 0)
			if request.Items != nil {
				seen := make(map[string]struct{})
				items := *request.Items
				for index, item := range items {
					item.ProductID = strings.TrimSpace(item.ProductID)
					if item.ProductID == "" {
						respondBadRequest(c, "productId is required")
						return
					}
					if item.Quantity <= 0 {
						respondBadRequest(c, "quantity must be greater than zero")
						return
					}
					if item.UnitPrice < 0 {
						respondBadRequest(c, "unitPrice must be zero or greater")
						return
					}
					items[index] = item
					if _, ok := seen[item.ProductID]; !ok {
						seen[item.ProductID] = struct{}{}
						productIDs = append(productIDs, item.ProductID)
					}
				}
				request.Items = &items
			}

			var purchase models.Purchase
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.First(&purchase, "id = ?", purchaseID).Error; err != nil {
					return err
				}

				if request.PurchaseDate != nil {
					parsed, err := time.ParseInLocation("2006-01-02", *request.PurchaseDate, time.Local)
					if err != nil {
						return err
					}
					purchase.PurchasedAt = parsed
				}

				if request.Items != nil {
					var count int64
					if err := tx.Model(&models.Product{}).Where("id in ?", productIDs).Count(&count).Error; err != nil {
						return err
					}
					if count != int64(len(productIDs)) {
						return errProductNotFound
					}

					if err := tx.Where("purchase_id = ?", purchaseID).Delete(&models.PurchaseItem{}).Error; err != nil {
						return err
					}

					total := 0.0
					for _, item := range *request.Items {
						total += item.Quantity * item.UnitPrice
						entry := models.PurchaseItem{
							PurchaseID: purchaseID,
							ProductID:  item.ProductID,
							Quantity:   item.Quantity,
							Price:      item.UnitPrice,
						}
						if err := tx.Omit("id").Create(&entry).Error; err != nil {
							return err
						}
					}
					purchase.Total = total
				}

				updates := map[string]interface{}{}
				if request.PurchaseDate != nil {
					updates["purchased_at"] = purchase.PurchasedAt
				}
				if request.Items != nil {
					updates["total"] = purchase.Total
				}
				if len(updates) > 0 {
					if err := tx.Model(&purchase).Updates(updates).Error; err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				if errors.Is(err, errProductNotFound) {
					respondNotFound(c, "product not found")
					return
				}
				if errors.Is(err, gorm.ErrRecordNotFound) {
					respondNotFound(c, "purchase not found")
					return
				}
				if strings.Contains(err.Error(), "cannot parse") || strings.Contains(err.Error(), "invalid syntax") {
					respondBadRequest(c, "purchaseDate must be in YYYY-MM-DD format")
					return
				}
				respondError(c, err)
				return
			}

			response, err := loadPurchaseResponse(purchaseID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					respondNotFound(c, "purchase not found")
					return
				}
				respondError(c, err)
				return
			}

			c.JSON(http.StatusOK, response)
		})

		r.DELETE("/purchases/:id", func(c *gin.Context) {
			purchaseID := strings.TrimSpace(c.Param("id"))
			if purchaseID == "" {
				respondBadRequest(c, "missing purchase id")
				return
			}

			result := db.Delete(&models.Purchase{}, "id = ?", purchaseID)
			if result.Error != nil {
				respondError(c, result.Error)
				return
			}
			if result.RowsAffected == 0 {
				respondNotFound(c, "purchase not found")
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	registerRoutes(router)
	registerRoutes(router.Group("/api/" + apiVersion))

	return router
}
