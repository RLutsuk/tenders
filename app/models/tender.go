package models

import "github.com/go-openapi/strfmt"

/*
"id": "550e8400-e29b-41d4-a716-446655440000",
"name": "Доставка товары Казань - Москва",
"description": "Нужно доставить оборудовоние для олимпиады по робототехники",
"status": "Created",
"serviceType": "Delivery",
"verstion": 1,
"createdAt": "2006-01-02T15:04:05Z07:00"*/

const (
	CREATEDTEN   string = "Created"
	PUBLISHEDTEN string = "Published"
	CLOSEDTEN    string = "Closed"
)

const (
	CONSTRUCTION string = "Construction"
	DELIVERY     string = "Delivery"
	MANUFACTURE  string = "Manufacture"
)

type Tender struct {
	Id              string          `json:"id,omitempty" gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	Name            string          `json:"name,omitempty" gorm:"column:name"`
	Description     string          `json:"description,omitempty" gorm:"column:description"`
	Status          string          `json:"status,omitempty" gorm:"column:status"`
	ServiceType     string          `json:"serviceType,omitempty" gorm:"column:service_type"`
	Version         uint32          `json:"version,omitempty" gorm:"column:version"`
	OrganizationId  string          `json:"organizationId,omitempty" gorm:"column:organization_id;type:uuid"`
	CreatorUsername string          `json:"creatorUsername,omitempty" gorm:"column:user_name"`
	Created         strfmt.DateTime `json:"createdAt,omitempty" gorm:"column:created_at"`
}
