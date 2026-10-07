package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	return db, nil
}
