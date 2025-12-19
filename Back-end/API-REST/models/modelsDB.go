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

type Backendtests struct {
	ID uint `gorm:"primaryKey"`

	Idgroup          uint
	Group            Groups `gorm:"foreignKey:Idgroup;references:ID"`
	Name             string
	Httptype         string
	Urlapi           string `json:"Urlapi"` 
	Requesttype      string
	Request          datatypes.JSON `gorm:"type:json"`
	Response         datatypes.JSON `gorm:"type:json"`
	ResponseHttpCode int
	Header           datatypes.JSON `gorm:"type:json"`
	Token            string
}

type Saveendpointresult struct {
	ID uint `gorm:"primaryKey"`

	Idgroup uint
	Group   Groups `gorm:"foreignKey:Idgroup;references:ID"`

	Idtest       uint
	Backendtests Backendtests `gorm:"foreignKey:Idtest;references:ID"`

	Testcasedescription string
	Testedinfrontend    bool
	Evidencefrontend    string
}
