package models

import "github.com/go-openapi/strfmt"

const (
	CREATEDBID   string = "Created"
	PUBLISHEDBID string = "Published"
	CANCELEDBID  string = "Canceled"
)

type Bid struct {
	Id          string          `json:"id,omitempty" gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	Name        string          `json:"name,omitempty" gorm:"column:name"`
	Description string          `json:"description,omitempty" gorm:"column:description"`
	Status      string          `json:"status,omitempty" gorm:"column:status"`
	TenderId    string          `json:"tenderId,omitempty" gorm:"column:tender_id;type:uuid"`
	AuthorType  string          `json:"authorType,omitempty" gorm:"column:author_type"`
	AuthorID    string          `json:"authorId,omitempty" gorm:"column:user_id;type:uuid"`
	Version     uint32          `json:"version,omitempty" gorm:"column:version"`
	Created     strfmt.DateTime `json:"createdAt,omitempty" gorm:"column:created_at"`
}

type Decision struct {
	Id       string          `gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	BidId    string          `gorm:"column:bid_id"`
	UserName string          `gorm:"column:user_name"`
	Decision string          `gorm:"column:decision"`
	Created  strfmt.DateTime `gorm:"column:created_at"`
}
