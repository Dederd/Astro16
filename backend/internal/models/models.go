package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ────────────────────────────────────────────────────────────
// JSON custom type helpers
// ────────────────────────────────────────────────────────────

type SelectedFlowersJSON []SelectedFlower

func (s SelectedFlowersJSON) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}
func (s *SelectedFlowersJSON) Scan(value interface{}) error {
	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into SelectedFlowersJSON", value)
	}
	return json.Unmarshal(raw, s)
}

type StringSliceJSON []string

func (s StringSliceJSON) Value() (driver.Value, error) {
	b, _ := json.Marshal(s)
	return string(b), nil
}
func (s *StringSliceJSON) Scan(value interface{}) error {
	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T", value)
	}
	return json.Unmarshal(raw, s)
}

// ────────────────────────────────────────────────────────────
// BouquetType
// ────────────────────────────────────────────────────────────

type BouquetType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Theme       string `json:"theme"`
}

// ────────────────────────────────────────────────────────────
// Flower — now stored in PostgreSQL
// ────────────────────────────────────────────────────────────

type FlowerDB struct {
	ID          string          `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Name        string          `gorm:"type:varchar(100);not null" json:"name"`
	NameID      string          `gorm:"type:varchar(100);not null" json:"name_id"`
	Description string          `gorm:"type:text" json:"description"`
	ImageURL    string          `gorm:"type:text" json:"image_url"`
	Emoji       string          `gorm:"type:varchar(10)" json:"emoji"`
	Price       int64           `gorm:"not null" json:"price"`
	IsAvailable bool            `gorm:"default:true" json:"is_available"`
	Stock       int             `gorm:"default:0" json:"stock"`
	Colors      StringSliceJSON `gorm:"type:jsonb" json:"colors"`
	Occasions   StringSliceJSON `gorm:"type:jsonb" json:"occasions"`
	Meaning     string          `gorm:"type:varchar(255)" json:"meaning"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (FlowerDB) TableName() string { return "flowers" }

func (f *FlowerDB) ToFlower() Flower {
	return Flower{
		ID:          f.ID,
		Name:        f.Name,
		NameID:      f.NameID,
		Description: f.Description,
		ImageURL:    f.ImageURL,
		Emoji:       f.Emoji,
		Price:       f.Price,
		IsAvailable: f.IsAvailable,
		Stock:       f.Stock,
		Colors:      []string(f.Colors),
		Occasions:   []string(f.Occasions),
		Meaning:     f.Meaning,
	}
}

type Flower struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	NameID      string   `json:"name_id"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image_url"`
	Emoji       string   `json:"emoji"`
	Price       int64    `json:"price"`
	IsAvailable bool     `json:"is_available"`
	Stock       int      `json:"stock"`
	Colors      []string `json:"colors"`
	Occasions   []string `json:"occasions"`
	Meaning     string   `json:"meaning"`
}

// ────────────────────────────────────────────────────────────
// Catalog Bouquet (pre-made designs from admin)
// ────────────────────────────────────────────────────────────

type CatalogBouquetDB struct {
	ID          string  `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Name        string  `gorm:"type:varchar(200);not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	ImageURL    string  `gorm:"type:text" json:"image_url"`
	Style       string  `gorm:"type:varchar(50)" json:"style"`
	Occasion    string  `gorm:"type:varchar(50)" json:"occasion"`
	Price       int64   `gorm:"not null" json:"price"`
	StemCount   int     `gorm:"default:0" json:"stem_count"`
	IsAvailable bool    `gorm:"default:true" json:"is_available"`
	Stock       int     `gorm:"default:0" json:"stock"`
	SortOrder   int     `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CatalogBouquetDB) TableName() string { return "catalog_bouquets" }

// ────────────────────────────────────────────────────────────
// SelectedFlower
// ────────────────────────────────────────────────────────────

type SelectedFlower struct {
	FlowerID string `json:"flower_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// ────────────────────────────────────────────────────────────
// AI Agent request/response types
// ────────────────────────────────────────────────────────────

type AgentVerifyRequest struct {
	BouquetTypeID string `json:"bouquet_type_id"`
	BouquetType   string `json:"bouquet_type"`
}

type AgentVerifyResponse struct {
	Message         string   `json:"message"`
	Recommendations []string `json:"recommendations"`
	Tips            string   `json:"tips"`
}

type GenerateBouquetRequest struct {
	BouquetTypeID   string           `json:"bouquet_type_id"`
	SelectedFlowers []SelectedFlower `json:"selected_flowers"`
	// Total stem count from user's actual selection — used to sync AI output
	TotalStemCount int `json:"total_stem_count"`
	// Optional hints from user
	StyleHint       string `json:"style_hint"`
	DescriptionHint string `json:"description_hint"`
}

type BouquetDesign struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Style       string      `json:"style"`
	ImagePrompt string      `json:"image_prompt"`
	SmallSize   SizeVariant `json:"small"`
	LargeSize   SizeVariant `json:"large"`
}

type SizeVariant struct {
	Label       string `json:"label"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
	StemCount   int    `json:"stem_count"`
}

type GenerateBouquetResponse struct {
	Designs []BouquetDesign `json:"designs"`
	Message string          `json:"message"`
}

// ────────────────────────────────────────────────────────────
// Order
// ────────────────────────────────────────────────────────────

type OrderDB struct {
	ID              string              `gorm:"primaryKey;type:varchar(30)" json:"id"`
	CustomerName    string              `gorm:"type:varchar(255);not null" json:"customer_name"`
	CustomerEmail   string              `gorm:"type:varchar(255);not null" json:"customer_email"`
	CustomerPhone   string              `gorm:"type:varchar(30);not null" json:"customer_phone"`
	BouquetTypeID   string              `gorm:"type:varchar(50)" json:"bouquet_type_id"`
	SelectedFlowers SelectedFlowersJSON `gorm:"type:jsonb;not null" json:"selected_flowers"`
	DesignID        string              `gorm:"type:varchar(50)" json:"design_id"`
	DesignName      string              `gorm:"type:varchar(255)" json:"design_name"`
	Size            string              `gorm:"type:varchar(20)" json:"size"`
	TotalAmount     int64               `gorm:"not null" json:"total_amount"`
	Status          string              `gorm:"type:varchar(30);default:'pending'" json:"status"`
	PaymentID       string              `gorm:"type:varchar(100)" json:"payment_id"`
	Notes           string              `gorm:"type:text" json:"notes"`
	// Shipping fields
	ShippingAddress  string `gorm:"type:text" json:"shipping_address"`
	ShippingCity     string `gorm:"type:varchar(100)" json:"shipping_city"`
	ShippingPostcode string `gorm:"type:varchar(10)" json:"shipping_postcode"`
	ShippingPhone    string `gorm:"type:varchar(30)" json:"shipping_phone"`
	// Courier fields
	CourierService   string `gorm:"type:varchar(30)" json:"courier_service"` // jne, jnt
	TrackingNumber   string `gorm:"type:varchar(100)" json:"tracking_number"`
	ShippingStatus   string `gorm:"type:varchar(50)" json:"shipping_status"`
	// Source: ai_generated or catalog
	OrderSource   string `gorm:"type:varchar(30);default:'ai_generated'" json:"order_source"`
	CatalogItemID string `gorm:"type:varchar(50)" json:"catalog_item_id"`
	// Auth: link to user account (nullable)
	UserID *uint `gorm:"index" json:"user_id,omitempty"`
	// Generate count for rate limiting
	GenerateCount int `gorm:"default:0" json:"generate_count"`
	// Biaya breakdown
	FlowerCost    int64 `gorm:"default:0" json:"flower_cost"`
	MakingFee     int64 `gorm:"default:5000" json:"making_fee"`
	AIFee         int64 `gorm:"default:0" json:"ai_fee"`
	ExtraQuotaFee int64 `gorm:"default:0" json:"extra_quota_fee"`
	ShippingCost  int64 `gorm:"default:0" json:"shipping_cost"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (OrderDB) TableName() string { return "orders" }

func (o *OrderDB) ToOrder() *Order {
	return &Order{
		ID:               o.ID,
		CustomerName:     o.CustomerName,
		CustomerEmail:    o.CustomerEmail,
		CustomerPhone:    o.CustomerPhone,
		BouquetTypeID:    o.BouquetTypeID,
		SelectedFlowers:  []SelectedFlower(o.SelectedFlowers),
		DesignID:         o.DesignID,
		DesignName:       o.DesignName,
		Size:             o.Size,
		TotalAmount:      o.TotalAmount,
		Status:           o.Status,
		PaymentID:        o.PaymentID,
		Notes:            o.Notes,
		ShippingAddress:  o.ShippingAddress,
		ShippingCity:     o.ShippingCity,
		ShippingPostcode: o.ShippingPostcode,
		ShippingPhone:    o.ShippingPhone,
		CourierService:   o.CourierService,
		TrackingNumber:   o.TrackingNumber,
		ShippingStatus:   o.ShippingStatus,
		OrderSource:      o.OrderSource,
		CatalogItemID:    o.CatalogItemID,
		UserID:           o.UserID,
		FlowerCost:       o.FlowerCost,
		MakingFee:        o.MakingFee,
		AIFee:            o.AIFee,
		ExtraQuotaFee:    o.ExtraQuotaFee,
		ShippingCost:     o.ShippingCost,
		CreatedAt:        o.CreatedAt,
	}
}

type Order struct {
	ID               string           `json:"id"`
	CustomerName     string           `json:"customer_name"`
	CustomerEmail    string           `json:"customer_email"`
	CustomerPhone    string           `json:"customer_phone"`
	BouquetTypeID    string           `json:"bouquet_type_id"`
	SelectedFlowers  []SelectedFlower `json:"selected_flowers"`
	DesignID         string           `json:"design_id"`
	DesignName       string           `json:"design_name"`
	Size             string           `json:"size"`
	TotalAmount      int64            `json:"total_amount"`
	Status           string           `json:"status"`
	PaymentID        string           `json:"payment_id"`
	Notes            string           `json:"notes"`
	ShippingAddress  string           `json:"shipping_address"`
	ShippingCity     string           `json:"shipping_city"`
	ShippingPostcode string           `json:"shipping_postcode"`
	ShippingPhone    string           `json:"shipping_phone"`
	CourierService   string           `json:"courier_service"`
	TrackingNumber   string           `json:"tracking_number"`
	ShippingStatus   string           `json:"shipping_status"`
	OrderSource      string           `json:"order_source"`
	CatalogItemID    string           `json:"catalog_item_id"`
	UserID           *uint            `json:"user_id,omitempty"`
	// Breakdown biaya
	FlowerCost    int64 `json:"flower_cost"`
	MakingFee     int64 `json:"making_fee"`
	AIFee         int64 `json:"ai_fee"`
	ExtraQuotaFee int64 `json:"extra_quota_fee"`
	ShippingCost  int64 `json:"shipping_cost"`
	CreatedAt     time.Time `json:"created_at"`
}

// ────────────────────────────────────────────────────────────
// Request types
// ────────────────────────────────────────────────────────────

type CreateOrderRequest struct {
	CustomerName     string           `json:"customer_name" binding:"required"`
	CustomerEmail    string           `json:"customer_email" binding:"required"`
	CustomerPhone    string           `json:"customer_phone" binding:"required"`
	BouquetTypeID    string           `json:"bouquet_type_id"`
	SelectedFlowers  []SelectedFlower `json:"selected_flowers"`
	DesignID         string           `json:"design_id"`
	DesignName       string           `json:"design_name"`
	Size             string           `json:"size"`
	TotalAmount      int64            `json:"total_amount" binding:"required"`
	Notes            string           `json:"notes"`
	ShippingAddress  string           `json:"shipping_address" binding:"required"`
	ShippingCity     string           `json:"shipping_city" binding:"required"`
	ShippingPostcode string           `json:"shipping_postcode" binding:"required"`
	ShippingPhone    string           `json:"shipping_phone" binding:"required"`
	CourierService   string           `json:"courier_service" binding:"required"`
	OrderSource      string           `json:"order_source"` // ai_generated | catalog
	CatalogItemID    string           `json:"catalog_item_id"`
	// Biaya breakdown
	FlowerCost     int64 `json:"flower_cost"`
	MakingFee      int64 `json:"making_fee"`
	AIFee          int64 `json:"ai_fee"`
	ExtraQuotaFee  int64 `json:"extra_quota_fee"`
	ShippingCost   int64 `json:"shipping_cost"`
}

type PaymentTokenRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

type PaymentTokenResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

// ────────────────────────────────────────────────────────────
// Admin types
// ────────────────────────────────────────────────────────────

type UpdateFlowerRequest struct {
	Price       *int64  `json:"price"`
	IsAvailable *bool   `json:"is_available"`
	Stock       *int    `json:"stock"`
	ImageURL    *string `json:"image_url"`
	Emoji       *string `json:"emoji"`
}

type UpdateOrderStatusRequest struct {
	Status         string `json:"status" binding:"required"`
	TrackingNumber string `json:"tracking_number"`
	ShippingStatus string `json:"shipping_status"`
	CourierService string `json:"courier_service"`
}

type CreateCatalogRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Style       string `json:"style"`
	Occasion    string `json:"occasion"`
	Price       int64  `json:"price" binding:"required"`
	StemCount   int    `json:"stem_count"`
	IsAvailable bool   `json:"is_available"`
	Stock       int    `json:"stock"`
	SortOrder   int    `json:"sort_order"`
}

// ────────────────────────────────────────────────────────────
// Generate session (rate limiting)
// ────────────────────────────────────────────────────────────

type GenerateSessionDB struct {
	SessionID     string    `gorm:"primaryKey;type:varchar(100)" json:"session_id"`
	GenerateCount int       `gorm:"default:0" json:"generate_count"`
	IsPaid        bool      `gorm:"default:false" json:"is_paid"`
	ExtraQuota    int       `gorm:"default:0" json:"extra_quota"`      // kuota tambahan yang dibeli
	ExtraQuotaFee int64     `gorm:"default:0" json:"extra_quota_fee"`  // total biaya kuota yang dibeli (Rp)
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (GenerateSessionDB) TableName() string { return "generate_sessions" }

// ────────────────────────────────────────────────────────────
// Design cache — simpan hasil generate berdasarkan hash kombinasi bunga
// ────────────────────────────────────────────────────────────

type DesignCacheDB struct {
	CacheKey    string    `gorm:"primaryKey;type:varchar(64)" json:"cache_key"` // SHA256 dari kombinasi bunga+type
	BouquetTypeID string  `gorm:"type:varchar(50)" json:"bouquet_type_id"`
	FlowerCombo string    `gorm:"type:text" json:"flower_combo"` // sorted JSON flower IDs+qty
	DesignsJSON string    `gorm:"type:jsonb" json:"designs_json"`
	Message     string    `gorm:"type:text" json:"message"`
	HitCount    int       `gorm:"default:0" json:"hit_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (DesignCacheDB) TableName() string { return "design_cache" }



type TrackingInfo struct {
	OrderID        string `json:"order_id"`
	TrackingNumber string `json:"tracking_number"`
	CourierService string `json:"courier_service"`
	ShippingStatus string `json:"shipping_status"`
	// From external courier API
	CourierData interface{} `json:"courier_data,omitempty"`
}

// ────────────────────────────────────────────────────────────
// User (Auth)
// ────────────────────────────────────────────────────────────

type UserDB struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                string    `gorm:"type:varchar(255);not null" json:"name"`
	Email               string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone               string    `gorm:"type:varchar(30)" json:"phone"`
	PasswordHash        string    `gorm:"type:varchar(255);not null" json:"-"`
	FreeGenerateCount   int       `gorm:"default:0" json:"free_generate_count"`   // berapa banyak free generate yang sudah dipakai
	ExtraQuota          int       `gorm:"default:0" json:"extra_quota"`           // kuota tambahan yang dibeli
	ExtraQuotaFee       int64     `gorm:"default:0" json:"extra_quota_fee"`       // total biaya kuota yang dibeli (Rp)
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (UserDB) TableName() string { return "users" }

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"user"`
}
