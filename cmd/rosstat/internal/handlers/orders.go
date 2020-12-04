package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AllOrders struct {
	SomeData string `json:"some_data"`
}

type ReadyToBuild struct {

}

func GetOrders(c echo.Context) error {
	ctx := c.(*RosContext)

	return ctx.JSON(http.StatusOK, AllOrders{
		SomeData: "foo",
	})
}

func GetSmallToBuildOrders(c echo.Context) error {
	ctx := c.(*RosContext)

	return ctx.JSON(http.StatusOK, AllOrders{
		SomeData: "foo",
	})
}

func PutSmallToBuildOrders(c echo.Context) error {
	ctx := c.(*RosContext)
	req := new(ReadyToBuild)

	if err := ctx.Bind(req); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, AllOrders{
		SomeData: "foo",
	})
}
