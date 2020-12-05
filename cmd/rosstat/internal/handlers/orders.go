package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"rosstat/cmd/rosstat/internal/db"
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

// GET /orders

func GetToBuildOrders(c echo.Context) error {
	ctx := c.(*RosContext)

	result, err := db.GetAllOrdersForCompletion()
	if err != nil{
		log.Println("error get all orders for completion: " + err.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// GET /orders/big/build - used by /orders/big page

func GetBigToBuildOrders(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	// remove it
	orderID = orderID

	// get the data by orderID
	result := []db.BigOrdersModel{
		{
			FormName: "Форма №1. Записная книжечка переписчика (бла бла бла балб лабла бал)",
			Total:    10,
			Built:    4,
		},
		{
			FormName: "Форма №1. Записная книжечка Котофея Матвеевича",
			Total:    0,
			Built:    0,
		},
		{
			FormName: "Форма №1. Записная книжечка Выгебало",
			Total:    0,
			Built:    0,
		},
		{
			FormName: "Форма №1. Записная книжечка кадавра",
			Total:    14,
			Built:    3,
		},
	}

	return ctx.JSON(http.StatusOK, result)
}
