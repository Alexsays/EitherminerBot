package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ConnectToDB allows user to connect to Postgres DB
func ConnectToDB(dbUser string, dbPassword string, dbName string) (*gorm.DB, error) {
	var connectionString = fmt.Sprintf(
		"host=postgres port=5432 sslmode=disable user=%s password=%s dbname=%s",
		dbUser, dbPassword, dbName,
	)

	return gorm.Open("postgres", connectionString)
}
