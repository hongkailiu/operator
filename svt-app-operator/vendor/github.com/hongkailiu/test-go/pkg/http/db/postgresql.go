package db

import (
	"github.com/jinzhu/gorm"
	// init postgres pkg
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// OpenPostgres opens a postgres db
func OpenPostgres(postgres string) (db *gorm.DB, err error) {
	return gorm.Open("postgres", postgres)
}
