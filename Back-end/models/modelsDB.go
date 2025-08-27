package models

import (
	"gorm.io/datatypes"
)

type Groups struct {
	ID         uint   `gorm:"primaryKey"`
	KeycloakID string `gorm:"column:keycloak_id"`
}

type Users struct {
	ID         uint   `gorm:"primaryKey"`
	KeycloakID string `gorm:"column:keycloak_id"`
}

type GroupsRelation struct { 
	ID uint `gorm:"primaryKey"`

	Iduser uint 
	User   Users `gorm:"foreignKey:Iduser;references:ID"`

	Idgroup uint 
	Group   Groups `gorm:"foreignKey:Idgroup;references:ID"`

	Accepted bool 
}

type BackendTests struct {
	ID uint `gorm:"primaryKey"`

	IdGroup uint 
	Group   Groups `gorm:"foreignKey:Idgroup;references:ID"`

	Httptype    string
	Urlapi      string
	Requesttype string
	Request     datatypes.JSON `gorm:"type:json"`
	Response    datatypes.JSON `gorm:"type:json"`
	Header      datatypes.JSON `gorm:"type:json"`
	Token       string
}

type SaveEndpointResult struct {
	ID uint `gorm:"primaryKey"`

	Idgroup uint
	Group   Groups `gorm:"foreignKey:Idgroup;references:ID"`

	Idtest uint
	BackendTests BackendTests `gorm:"foreignKey:Idtest;references:ID"`

	Testcasedescription string
	TestedInfrontend    bool
	Evidencefrontend    string
}
