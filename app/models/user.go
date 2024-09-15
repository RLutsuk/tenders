package models

import "github.com/go-openapi/strfmt"

type User struct {
	Id         string
	Username   string
	First_name string
	Last_name  string
	Created    strfmt.DateTime
	Updated    strfmt.DateTime
}
