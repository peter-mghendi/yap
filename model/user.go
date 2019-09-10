package model

import (
	"errors"
	"net/http"

	"github.com/l3njo/yap-api/db"

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
	Role      UserRole   `json:"role"`
	Posts     []Post     `json:"posts,omitempty" sql:"-" gorm:"foreignkey:Creator"`
	Reactions []Reaction `json:"reactions,omitempty" sql:"-" gorm:"foreignkey:User"`
}

// UserRole represents a user rank
type UserRole string

// UserRoles represent various user ranks
const (
	UserReader UserRole = "reader"
	UserEditor UserRole = "editor"
	UserKeeper UserRole = "keeper"
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
	if db.DB.Model(&User{}).Count(&count); count == 0 {
		u.Role = UserKeeper
	}

	if err = db.DB.Create(u).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

// Read fetches a User
func (u *User) Read() (int, error) {
	if err := db.DB.Set("gorm:auto_preload", true).First(u).Error; gorm.IsRecordNotFoundError(err) {
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

	if err := db.DB.Model(u).Updates(user).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := db.DB.First(u).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes a User
func (u *User) Delete() (int, error) {
	db := db.DB.Delete(u)
	if num, err := db.RowsAffected, db.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// ValidateAuth checks user details format
func (u *User) ValidateAuth() (int, error) {
	code := http.StatusOK
	if u.Mail == "" || u.Pass == "" {
		code = http.StatusBadRequest
		return code, errors.New(http.StatusText(code))
	}

	return code, errors.New(http.StatusText(code))
}

// TryAuth checks user credentials
func (u *User) TryAuth() (int, error) {
	pass := []byte(u.Pass)
	user := &User{Mail: u.Mail}
	if err := db.DB.Find(user).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Pass), pass)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	*u = *user
	return http.StatusAccepted, nil
}

// ReadAllUsers fetches all Users
func ReadAllUsers() ([]User, int, error) {
	users := []User{}
	if err := db.DB.Set("gorm:auto_preload", true).Find(&users).Error; gorm.IsRecordNotFoundError(err) {
		return users, http.StatusNotFound, err
	}

	return users, http.StatusOK, nil
}

// CountUsers counts specified type of users
func CountUsers(u *User) (int, int, error) {
	var count int
	users := []User{}
	if err := db.DB.Where(u).Find(&users).Count(&count).Error; err != nil {
		return count, http.StatusOK, err
	}

	return count, http.StatusOK, nil
}
