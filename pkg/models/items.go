package models

import "github.com/google/uuid"

type Items struct {
	BaseModel
	ItemName string     `json:"item_name"`
	UserID   *uuid.UUID `json:"user_id"`
}
