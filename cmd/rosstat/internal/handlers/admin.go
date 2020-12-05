package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"rosstat/cmd/rosstat/internal/db"
)

func GetUsers(c echo.Context) error {
	ctx := c.(*RosContext)

	users, err := db.GetPublicUsers()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, users)
}

func AddUser(c echo.Context) error {
	ctx := c.(*RosContext)
	user := new(db.User)
	if err := ctx.Bind(user); err != nil {
		return err
	}

	if err := db.AddUser(user.Login, user.Password, user.Role); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func DeleteUser(c echo.Context) error {
	ctx := c.(*RosContext)
	login := c.Param("login")

	if err := db.DeleteUser(login); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
