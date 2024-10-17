package dao

import (
	"fastgin/database"
	"gorm.io/gorm"
)

type DictionaryDAO struct {
	db *gorm.DB
}

func NewDictionaryDAO() *DictionaryDAO {
	return &DictionaryDAO{db: database.DB}
}
