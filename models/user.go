package models

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User is a registered user
type User struct {
	Base
	Name      string     `json:"name"`
	Mail      string     `json:"mail"`
	Pass      string     `json:"pass"`
	Auth      string     `json:"auth"`
	Life      string     `json:"life"`
	Role      userRole   `json:"role"`
	Posts     []Post     `json:"posts,omitempty" sql:"-" gorm:"foreignkey:Creator"`
	Reactions []Reaction `json:"reactions,omitempty" sql:"-" gorm:"foreignkey:User"`
}

type userRole int

// userRoles represent various user ranks
const (
	UserReader userRole = iota + 1
	UserEditor
	UserKeeper
)

// Create makes a User
// First user is automatically promoted to "UserKeeper" role
func (u *User) Create() (int, error) {
	var count int
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	u.Pass = string(hash)
	u.Role = UserReader
	if DB.Model(&User{}).Count(&count); count == 0 {
		u.Role = UserKeeper
	}

	if err = DB.Create(u).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

// Read fetches a User
func (u *User) Read() (int, error) {
	if err := DB.Set("gorm:auto_preload", true).First(u).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Update edits a User
func (u *User) Update() (int, error) {
	if u.Pass != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		u.Pass = string(hash)
	}

	user := User{
		Name: u.Name,
		Mail: u.Mail,
		Pass: u.Pass,
		Role: u.Role,
		Life: u.Life,
	}

	if err := DB.Model(u).Updates(user).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes a User
func (u *User) Delete() (int, error) {
	db := DB.Delete(u)
	if num, err := db.RowsAffected, db.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusAccepted, nil
}
