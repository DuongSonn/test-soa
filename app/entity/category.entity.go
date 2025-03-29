package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sondth-test_soa/utils"
)

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"varchar(255);not null"`
	NameSlug    string    `json:"-" gorm:"varchar(255);not null"`
	Description *string   `json:"description" gorm:"text"`
	CreatedAt   int64     `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at,omitempty" gorm:"autoUpdateTime:milli"`

	// Response fields
	ProductCount int64 `json:"product_count" gorm:"-"`
}

func NewCategory() *Category {
	return &Category{
		ID:        uuid.New(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (Category) TableName() string {
	return "categories"
}

func (e *Category) BeforeSave(tx *gorm.DB) error {
	e.UpdatedAt = time.Now().Unix()
	if e.Name != "" {
		e.NameSlug = utils.ConvertToSlug(e.Name)
	}

	return nil
}
