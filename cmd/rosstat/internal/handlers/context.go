package handlers

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

type RosContext struct {
	echo.Context
	db *sql.DB
}

func (r RosContext) DB() *sql.DB {
	return r.db
}

func RosContextMiddleware(dbConnStr string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			db, err := sql.Open("postgres", dbConnStr)
			if err != nil {
				return fmt.Errorf("unable to connect to db: %s", err)
			}
			defer func() {
				if err = db.Close(); err != nil {
					c.Logger().Error(err)
				}
			}()

			cc := &RosContext{
				Context: c,
				db: db,
			}
			return next(cc)
		}
	}
}
