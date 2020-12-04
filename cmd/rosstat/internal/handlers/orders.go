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

// GET /orders

func GetToBuildOrders(c echo.Context) error {
	ctx := c.(*RosContext)

	result := []db.OrdersModel{
		{
			ID:            1,
			Num:           1,
			OrderCaption:  "О-20-123-РОССТАТ 2",
			Customer:      "Росстат",
			Address:       "107123, Москва",
			Run:           270,
			AmountPallets: 1,
			AmountBoxes:   1,
			SubOrders: []db.SubOrderModel{
				{
					IsSmall:       false,
					OrderCaption:  "О-20-123-РОССТАТ 2 (1-26)",
					AmountPallets: 1,
					AmountBoxes:   1,
				},
				{
					IsSmall:       true,
					OrderCaption:  "О-20-123-РОССТАТ 2 (27)",
					AmountPallets: 1,
					AmountBoxes:   1,
				},
			},
		},
		{
			ID:            2,
			Num:           2,
			OrderCaption:  "О-22-355-РОССТАТ 1",
			Customer:      "Росстат",
			Address:       "107123, Москва",
			Run:           1650,
			AmountPallets: 0,
			AmountBoxes:   0,
			SubOrders: []db.SubOrderModel{
				{
					IsSmall:       false,
					OrderCaption:  "О-22-355-РОССТАТ 1 (1-26)",
					AmountPallets: 0,
					AmountBoxes:   0,
				},
				{
					IsSmall:       true,
					OrderCaption:  "О-22-355-РОССТАТ 1 (27)",
					AmountPallets: 0,
					AmountBoxes:   0,
				},
			},
		},
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

// GET /orders/big/build/:id from page /orders/small/:id

func GetSmallToBuildOrders(c echo.Context) error {
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

func FinishSmallToBuildOrders(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	req := new(db.FinishSmallOrderModel)
	if err := ctx.Bind(req); err != nil {
		return err
	}

	log.Println(orderID)
	log.Println(req.Boxes)

	return ctx.NoContent(http.StatusNoContent)
}
