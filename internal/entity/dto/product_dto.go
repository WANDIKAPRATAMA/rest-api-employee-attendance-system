package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string    `json:"name" validate:"required"`
	SKU         string    `json:"sku" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Description string    `json:"description"`
}

type UpdateProductRequest struct {
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	CategoryID  uuid.UUID `json:"category_id"`
	Description string    `json:"description"`
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	CategoryID  uuid.UUID `json:"category_id"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type CreateProductCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateProductCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
}

type CreateProductStockRequest struct {
	ProductID           uuid.UUID `json:"product_id" validate:"required"`
	WarehouseLocationID uuid.UUID `json:"warehouse_location_id" validate:"required"`
	Quantity            int       `json:"quantity" validate:"required,min=0"`
}

type UpdateProductStockRequest struct {
	Quantity int `json:"quantity" validate:"min=0"`
}

type ProductStockResponse struct {
	ID                  uuid.UUID `json:"id"`
	ProductID           uuid.UUID `json:"product_id"`
	WarehouseLocationID uuid.UUID `json:"warehouse_location_id"`
	Quantity            int       `json:"quantity"`
	Status              string    `json:"status"`
	UpdatedAt           string    `json:"updated_at"`
}

type ApiResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload"`
}

// Lists product response
type ProductListResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	SKU          string    `json:"sku"`
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// WarehouseLocationListResponse
type WarehouseLocationListResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateWarehouseLocationRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// UpdateWarehouseLocationRequest
type UpdateWarehouseLocationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// WarehouseLocationResponse
type WarehouseLocationResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductStockListResponse: Enriched dengan ProductName dan WarehouseName
type ProductStockListResponse struct {
	ID                  uuid.UUID `json:"id"`
	ProductID           uuid.UUID `json:"product_id"`
	ProductName         string    `json:"product_name"`
	WarehouseLocationID uuid.UUID `json:"warehouse_location_id"`
	WarehouseName       string    `json:"warehouse_name"`
	Quantity            int       `json:"quantity"`
	Status              string    `json:"status"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// LowStockDetail untuk item low-stock dengan detail
type LowStockDetail struct {
	ProductID      uuid.UUID `json:"product_id"`
	ProductName    string    `json:"product_name"`
	WarehouseID    uuid.UUID `json:"warehouse_id"`
	WarehouseName  string    `json:"warehouse_name"`
	Quantity       int       `json:"quantity"`
	Status         string    `json:"status"`
	UpdatedByEmail string    `json:"updated_by_email"`
	UpdatedByName  string    `json:"updated_by_name"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// OutOfStockDetail serupa dengan low-stock
type OutOfStockDetail struct {
	ProductID      uuid.UUID `json:"product_id"`
	ProductName    string    `json:"product_name"`
	WarehouseID    uuid.UUID `json:"warehouse_id"`
	WarehouseName  string    `json:"warehouse_name"`
	Quantity       int       `json:"quantity"`
	Status         string    `json:"status"`
	UpdatedByEmail string    `json:"updated_by_email"`
	UpdatedByName  string    `json:"updated_by_name"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RecentAddition untuk produk baru ditambahkan
type RecentAddition struct {
	ProductID      uuid.UUID `json:"product_id"`
	ProductName    string    `json:"product_name"`
	CreatedByEmail string    `json:"created_by_email"`
	CreatedByName  string    `json:"created_by_name"`
	CreatedAt      time.Time `json:"created_at"`
}

// DashboardResponse yang lebih modular dan kompleks
type DashboardResponse struct {
	TotalStock       int64              `json:"total_stock"`
	NumberOfProducts int64              `json:"number_of_products"`
	LowStockItems    []LowStockDetail   `json:"low_stock_items"`    // List detail low-stock
	OutOfStockItems  []OutOfStockDetail `json:"out_of_stock_items"` // List detail out-of-stock
	RecentAdditions  []RecentAddition   `json:"recent_additions"`   // List produk baru (misal last 5)
}

// PaginationRequest untuk query param
type PaginationRequest struct {
	Page   int    `query:"page" validate:"min=1"`
	Limit  int    `query:"limit" validate:"min=1,max=100"`
	Search string `query:"search"`
	SortBy string `query:"sort_by" validate:"oneof=name sku created_at"`
	Order  string `query:"order" validate:"oneof=asc desc"`
	// Filter spesifik
	CategoryID uuid.UUID `query:"category_id"`
	Status     string    `query:"status" validate:"oneof=available low-stock out-of-stock"`
}

type Pagination struct {
	HasNextPage bool `json:"has_next_page"`
	NextPage    *int `json:"next_page"`
	CurrentPage int  `json:"current_page"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
}

type ProductCategoryListResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
