package handler

import (
	"net/http"

	"github.com/l3njo/yap/model"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// UserResponse is a response containing one User
type UserResponse struct {
	Response
	model.User `json:"data"`
}

// UsersResponse is a response containing a slice of Users
type UsersResponse struct {
	Response
	Users []model.User `json:"data"`
}

// GetUsers handles the "/users" route.
func GetUsers(c echo.Context) error {
	resp, status := UsersResponse{}, 0
	users, status, err := model.ReadAllUsers()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	for _, user := range users {
		user.Pass = ""
		resp.Users = append(resp.Users, user)
	}

	resp.Status, resp.Message = true, http.StatusText(status)
	return c.JSON(status, resp)
}

// GetUserByID handles the "/users/:id" route.
func GetUserByID(c echo.Context) error {
	resp, status := UserResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user := model.User{
		Base: model.Base{ID: id},
	}

	status, err := user.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Pass = ""
	resp.Status, resp.Message, resp.User = true, http.StatusText(status), user
	return c.JSON(status, resp)
}

// UpdateUser handles the "/users/:id/update" route.
// TODO Check if user in context matches user, else get user from context
func UpdateUser(c echo.Context) error {
	resp, status := UserResponse{}, 0
	user, u := model.User{}, model.User{}
	if err := c.Bind(&u); err != nil {
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

	if status, err := user.Read(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Name, user.Mail, user.Life = u.Name, u.Mail, u.Life
	status, err := user.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Pass = ""
	resp.Status, resp.Message, resp.User = true, http.StatusText(status), user
	return c.JSON(status, resp)
}

// AssignUser handles the "/users/:id/assign" route.
// TODO Check if user in context is keeper
func AssignUser(c echo.Context) error {
	resp, status := UserResponse{}, 0
	user, u := model.User{}, model.User{}
	if err := c.Bind(&u); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if u.Role != model.UserKeeper &&
		u.Role != model.UserEditor &&
		u.Role != model.UserReader {
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

	if status, err := user.Read(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if u.Role != model.UserKeeper && user.Role == model.UserKeeper {
		count, status, err := model.CountUsers(&model.User{Role: model.UserKeeper})
		if err != nil {
			status = http.StatusInternalServerError
			resp.Message = http.StatusText(status)
			return c.JSON(status, resp)
		}

		if count == 1 {
			status = http.StatusNotModified
			resp.Message = "Sole keeper"
			return c.JSON(status, resp)
		}
	}

	user.Role = u.Role
	status, err := user.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user.Pass = ""
	resp.Status, resp.Message, resp.User = true, http.StatusText(status), user
	return c.JSON(status, resp)
}

// DeleteUser handles the "/users/:id/delete" route.
// TODO Check if user in context matches user, else get user from context
func DeleteUser(c echo.Context) error {
	resp, status := UserResponse{}, 0
	user := model.User{
		Base: model.Base{ID: uuid.FromStringOrNil(c.Param("id"))},
	}

	status, err := user.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if user.Role == model.UserKeeper {
		count, status, err := model.CountUsers(&model.User{Role: model.UserKeeper})
		if err != nil {
			status = http.StatusInternalServerError
			resp.Message = http.StatusText(status)
			return c.JSON(status, resp)
		}

		if count == 1 {
			status = http.StatusNotModified
			resp.Message = "Sole keeper"
			return c.JSON(status, resp)
		}
	}

	status, err = user.Delete()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message = true, http.StatusText(status)
	return c.JSON(status, resp)
}
