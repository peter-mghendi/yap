package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/xo/dburl"
)

// DB is the database
var DB *gorm.DB

// parseDBURL returns parses a database url into individual components
func parseDBURL(url string) (string, string, error) {
	u, err := dburl.Parse(url)
	if err != nil {
		return "", "", err
	}

	return u.Driver, fmt.Sprintf("%s sslmode=disable", u.DSN), nil
}

// InitDB sets up the databases
func InitDB(url string) error {
	dialect, uri, err := parseDBURL(url)
	if err != nil {
		return err
	}

	conn, err := gorm.Open(dialect, uri)
	if err != nil {
		return err
	}

	DB = conn
	DB.Debug().AutoMigrate(&User{}, &Article{}, &Gallery{}, &Flicker{}, &Reaction{})
	return nil
}
