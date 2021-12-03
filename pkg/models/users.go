package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	BaseModel
	Username string `json:"username" gorm:"unique"`
	FullName string `json:"full_name"`
	Password string `json:"-"`

	// Foreign Keys
	Items []Items `json:"items"`

	//
}

func (u *Users) BeforeCreate(_ *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
