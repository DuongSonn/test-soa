package entity

import (
	"sondth-test_soa/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username  string    `json:"username" gorm:"varchar(255);not null;unique"`
	Password  string    `json:"password" gorm:"varchar(255);not null"`
	Fullname  string    `json:"fullname" gorm:"varchar(255);not null"`
	Role      string    `json:"role" gorm:"varchar(255);not null"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func NewUser() *User {
	return &User{
		ID:        uuid.New(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().Unix()
	if u.Password != "" {
		u.Password, err = utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}

	return
}

func (u *User) CheckPassword(password string) error {
	return utils.CheckPasswordHash(password, u.Password)
}

func (u *User) IsAdmin() bool {
	return u.Role == ROLE_ADMIN
}
