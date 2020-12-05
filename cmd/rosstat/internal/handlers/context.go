package handlers

import (
	"database/sql"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"rosstat/cmd/rosstat/internal/db"
)

type RosContext struct {
	echo.Context
	db   *sql.DB
	user db.PublicUser
}

func (r RosContext) DB() *sql.DB {
	return r.db
}

func (r RosContext) User() db.PublicUser {
	return r.user
}

func RosContextMiddleware(dbConnStr string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			dbConn, err := sql.Open("postgres", dbConnStr)
			if err != nil {
				return fmt.Errorf("unable to connect to db: %s", err)
			}
			defer func() {
				if err = dbConn.Close(); err != nil {
					c.Logger().Error(err)
				}
			}()

			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			login := claims["login"].(string)
			role := claims["role"].(string)

			cc := &RosContext{
				Context: c,
				db:      dbConn,
				user: db.PublicUser{
					Login: login,
					Role:  role,
				},
			}
			return next(cc)
		}
	}
}
