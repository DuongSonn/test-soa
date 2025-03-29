package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Rating    float64   `json:"rating" gorm:"type:decimal(10,2);not null"`
	Comment   string    `json:"comment" gorm:"text"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`

	// Relations
	Product Product `json:"product"`
	User    User    `json:"user"`
}

func NewReview() *Review {
	return &Review{
		ID:        uuid.New(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (Review) TableName() string {
	return "reviews"
}

func (e *Review) BeforeSave(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now().Unix()
	return
}
