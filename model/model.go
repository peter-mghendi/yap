package model

import (
	"github.com/l3njo/yap/db"
)

// InitDB connects to and sets up the database
func InitDB(url string) error {
	if err := db.Init(url); err != nil {
		return err
	}
	if err := db.DB.Debug().AutoMigrate(&User{}, &Article{}, &Gallery{}, &Flicker{}, &Question{}, &Response{}, &Reaction{}).Error; err != nil {
		return err
	}

	return nil
}
