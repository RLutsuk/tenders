package models

import "github.com/go-openapi/strfmt"

type User struct {
	Id         string          `gorm:"column:id"`
	Username   string          `gorm:"column:username"`
	First_name string          `gorm:"column:first_name"`
	Last_name  string          `gorm:"column:last_name"`
	Created    strfmt.DateTime `gorm:"column:created_at"`
	Updated    strfmt.DateTime `gorm:"column:updated_at"`
}
