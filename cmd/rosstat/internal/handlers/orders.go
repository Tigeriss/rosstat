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

// GET /orders

type SubOrderModel struct {
	IsSmall       bool   `json:"is_small"`      // Тип подзаказа. true - small. false - big
	OrderCaption  string `json:"order_caption"` // Заказ
	AmountPallets int    `json:"amount_pallets"`
	AmountBoxes   int    `json:"amount_boxes"`
}

type OrdersModel struct {
	ID            int             `json:"id"`
	Num           int             `json:"num"`           // Номер
	OrderCaption  string          `json:"order_caption"` // Заказ
	Customer      string          `json:"customer"`      // Заказчик
	Address       string          `json:"address"`       // Адрес
	Run           int             `json:"run"`           // Тираж
	AmountPallets int             `json:"amount_pallets"`
	AmountBoxes   int             `json:"amount_boxes"`
	SubOrders     []SubOrderModel `json:"sub_orders"`
}

func GetToBuildOrders(c echo.Context) error {
	ctx := c.(*RosContext)

	result := []OrdersModel{
		{
			ID:            1,
			Num:           1,
			OrderCaption:  "О-20-123-РОССТАТ 2",
			Customer:      "Росстат",
			Address:       "107123, Москва",
			Run:           270,
			AmountPallets: 1,
			AmountBoxes:   1,
			SubOrders: []SubOrderModel{
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
			SubOrders: []SubOrderModel{
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
