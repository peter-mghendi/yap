package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/l3njo/yap/model"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	User uuid.UUID      `json:"user"`
	Role model.UserRole `json:"role"`
	jwt.StandardClaims
}

func createAuthString(user model.User) (string, error) {
	claims := &JwtCustomClaims{
		User: user.ID,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	authString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return authString, nil
}

// JoinUser handles the "/users/join" route.
func JoinUser(c echo.Context) error {
	resp, status := UserResponse{}, 0
	user, u := model.User{}, model.User{}
	if err := c.Bind(&u); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Name, user.Mail, user.Pass, user.Life = u.Name, u.Mail, u.Pass, u.Life
	status, err := user.Create()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Pass = ""
	authString, err := createAuthString(user)
	if err != nil {
		status = http.StatusInternalServerError
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Auth = authString
	resp.Status, resp.Message, resp.User = true, http.StatusText(status), user
	return c.JSON(status, resp)
}

// AuthUser handles the "users/auth" route.
func AuthUser(c echo.Context) error {
	resp, status := UserResponse{}, 0
	user := model.User{}
	err := c.Bind(&user)
	if err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err = user.ValidateAuth()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err = user.TryAuth()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Pass = ""
	authString, err := createAuthString(user)
	if err != nil {
		status = http.StatusInternalServerError
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Auth = authString
	resp.Status, resp.Message, resp.User = true, http.StatusText(status), user
	return c.JSON(status, resp)
}

// UpdatePass handles the "/users/:id/change" route.
func UpdatePass(c echo.Context) error {
	resp, status := UserResponse{}, 0
	user, u := model.User{}, map[string]string{}
	if err := c.Bind(&u); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if u["current"] == "" || u["updated"] == "" {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.ID = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(user.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Pass = u["current"]
	status, err := user.ValidateAuth()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err = user.TryAuth()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Pass = u["updated"]
	status, err = user.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Pass = ""
	resp.Status, resp.Message, resp.User = true, http.StatusText(status), user
	return c.JSON(status, resp)
}
