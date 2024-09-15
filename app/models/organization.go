package models

import (
	"github.com/go-openapi/strfmt"
)

const (
	IE  string = "IE"
	LLC string = "LLC"
	JSC string = "JSC"
)

type Organization struct {
	Id          string          `json:"id,omitempty" gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	Name        string          `json:"name,omitempty" gorm:"column:name"`
	Description string          `json:"description,omitempty" gorm:"column:description"`
	Type        string          `json:"type,omitempty" gorm:"column:type"`
	Created     strfmt.DateTime `json:"created,omitempty" gorm:"column:created_at"`
	Updated     strfmt.DateTime `json:"updated,omitempty" gorm:"column:updated_at"`
}

type OrganizationResponsible struct {
	Id           string `gorm:"id"`
	Organization string `gorm:"organization_id"`
	User         string `gorm:"user_id"`
}
