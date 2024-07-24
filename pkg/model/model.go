package models

import "gorm.io/gorm"

type TableName struct {
	Id     uint    `json:"id"`
	Name   string  `json:"name,omitempty"`
	Fields []Field `json:"fields,omitempty"`
}
type Field struct {
	gorm.Model
	Name    string `json:"name,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Type    string `json:"type,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type Modeltor interface {
	CreateTable() (table interface{}, err error)

	DropTable(id uint) error
}
