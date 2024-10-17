package database

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

// GetDatabaseNames retrieves the list of database names.
func GetDatabaseNames(db *gorm.DB) ([]string, error) {
	var databases []string
	query := "SHOW DATABASES"

	if err := db.Raw(query).Scan(&databases).Error; err != nil {
		log.Printf("Error fetching database names: %v", err)
		return nil, err
	}
	return databases, nil
}

// GetTableNames retrieves the list of table names from the specified database.
func GetTableNames(db *gorm.DB, databaseName string) ([]string, error) {
	var tables []string
	query := fmt.Sprintf("SHOW TABLES FROM %s", databaseName)

	if err := db.Raw(query).Scan(&tables).Error; err != nil {
		log.Printf("Error fetching table names from database %s: %v", databaseName, err)
		return nil, err
	}
	return tables, nil
}

// ColumnInfo represents the information of a table column.
type ColumnInfo struct {
	ColumnName    string `json:"column_name"`
	ColumnType    string `json:"column_type"`
	IsNullable    string `json:"is_nullable"`
	ColumnKey     string `json:"column_key"`
	ColumnDefault string `json:"column_default"`
	Extra         string `json:"extra"`
	ColumnComment string `json:"column_comment"`
}

func (c ColumnInfo) IsPriKey() bool {
	return c.ColumnKey == "PRI"
}

// GetTableInfo retrieves the information of the specified table from the specified database, including column comments.
func GetTableInfo(db *gorm.DB, databaseName, tableName string) ([]ColumnInfo, error) {
	var tableInfo []ColumnInfo
	query := fmt.Sprintf("SELECT column_name, column_type, is_nullable, column_key, column_default, extra, column_comment FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'", databaseName, tableName)
	if err := db.Raw(query).Scan(&tableInfo).Error; err != nil {
		log.Printf("Error fetching table info from table %s in database %s: %v", tableName, databaseName, err)
		return nil, err
	}
	return tableInfo, nil
}
