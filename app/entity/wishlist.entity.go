package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wishlist struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`

	// Relations
	User    User    `json:"user"`
	Product Product `json:"product"`
}

func NewWishlist() *Wishlist {
	return &Wishlist{
		ID:        uuid.New(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (Wishlist) TableName() string {
	return "wishlists"
}

func (e *Wishlist) BeforeSave(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now().Unix()
	return
}
