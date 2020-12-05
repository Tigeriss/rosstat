package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"rosstat/cmd/rosstat/internal/db"
)

func GetReadyForShipment(c echo.Context) error {
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
		},
		{
			ID:            2,
			Num:           2,
			OrderCaption:  "О-22-355-РОССТАТ 1",
			Customer:      "Росстат",
			Address:       "107123, Москва",
			Run:           1650,
			AmountPallets: 1,
			AmountBoxes:   8,
		},
	}

	return ctx.JSON(http.StatusOK, result)
}

func GetPalletShipment(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	orderID = orderID

	result := []db.ShipmentPalletModel{
		{
			Num:         1,
			PalletNum:   1,
			Barcode:     "111111111",
			AmountBoxes: 96,
		},
		{
			Num:         2,
			PalletNum:   2,
			Barcode:     "111111112",
			AmountBoxes: 96,
		},
		{
			Num:         3,
			PalletNum:   3,
			Barcode:     "111111113",
			AmountBoxes: 96,
		},
		{
			Num:         4,
			PalletNum:   4,
			Barcode:     "111111114",
			AmountBoxes: 96,
		},
		{
			Num:         5,
			PalletNum:   5,
			Barcode:     "111111115",
			AmountBoxes: 96,
		},
		{
			Num:         6,
			PalletNum:   6,
			Barcode:     "111111116",
			AmountBoxes: 96,
		},
		{
			Num:         7,
			PalletNum:   7,
			Barcode:     "111111117",
			AmountBoxes: 96,
		},
	}

	return ctx.JSON(http.StatusOK, result)
}

func FinishPalletShipment(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	log.Println(orderID)

	return ctx.NoContent(http.StatusNoContent)
}
