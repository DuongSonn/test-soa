package entity

import (
	"sondth-test_soa/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	PRODUCT_STATUS_IN_STOCK     = "in_stock"
	PRODUCT_STATUS_OUT_OF_STOCK = "out_of_stock"
)

type Product struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"varchar(255);not null"`
	NameSlug    string    `json:"-" gorm:"varchar(255);not null"`
	Description *string   `json:"description" gorm:"text"`
	Image       *string   `json:"image" gorm:"varchar(255)"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity    uint64    `json:"quantity" gorm:"type:bigint unsigned;not null"`
	CreatedAt   int64     `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at,omitempty" gorm:"autoUpdateTime:milli"`

	// Relations
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	Category   Category  `json:"category"`

	// Response fields
	Status string `json:"status" gorm:"-:all"`
}

func NewProduct() *Product {
	return &Product{
		ID:        uuid.New(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (Product) TableName() string {
	return "products"
}

func (e *Product) BeforeSave(tx *gorm.DB) error {
	e.UpdatedAt = time.Now().Unix()
	if e.Name != "" {
		e.NameSlug = utils.ConvertToSlug(e.Name)
	}

	return nil
}
