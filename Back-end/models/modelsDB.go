package models

import (
	"gorm.io/datatypes"
)

type Groups struct {
	ID         uint `gorm:"primaryKey"`
	KeycloakID string
}

type Users struct {
	ID         uint `gorm:"primaryKey"`
	KeycloakID string
}

type Groups_relation struct {
	ID uint `gorm:"primaryKey"`

	idUser uint
	User   Users `gorm:"foreignKey:idUser;references:ID"`

	idGroup uint
	Group   Groups `gorm:"foreignKey:idGroup;references:ID"`
}

type Backend_tests struct {
	ID uint `gorm:"primaryKey"`

	idGroup uint
	Group   Groups `gorm:"foreignKey:idGroup;references:ID"`

	http_type    string
	url_api      string
	request_type string
	request      datatypes.JSON `gorm:"type:json"`
	response     datatypes.JSON `gorm:"type:json"`
	header       datatypes.JSON `gorm:"type:json"`
	token        string
}

type Save_endpoint_result struct {
	ID uint `gorm:"primaryKey"`

	idGroup uint
	Group   Groups `gorm:"foreignKey:idGroup;references:ID"`

	idTest uint
	Backend_tests   Backend_tests `gorm:"foreignKey:idTest;references:ID"`

	test_case_description    string
	tested_in_frontend      bool
	Evidence_frontend string
}
