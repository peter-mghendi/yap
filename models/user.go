package models

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User is a registered user
type User struct {
	Base
	Name string   `json:"name,omitempty"`
	Mail string   `json:"mail,omitempty"`
	Pass string   `json:"pass,omitempty"`
	Auth string   `json:"auth,omitempty"`
	Life string   `json:"life,omitempty"`
	Role userRole `json:"role,omitempty"`
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
	if err := DB.First(u).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
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

// ReadReactions fetches a User's Reactions
func (u *User) ReadReactions(reaction reactionType) ([]Reaction, int, error) {
	return ReadReactions(&Reaction{User: u.ID, Type: reaction})
}
