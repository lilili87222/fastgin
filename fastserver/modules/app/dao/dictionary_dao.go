package dao

import (
"fastgin/database"
"gorm.io/gorm"
 _ "fastgin/modules/app/model"
)

type DictionaryDAO struct {
 db *gorm.DB
}

func NewDictionaryDAO() *DictionaryDAO {
 return &DictionaryDAO{db: database.DB}
}
