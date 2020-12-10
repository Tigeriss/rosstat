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

	tx, err := ctx.DB().Begin()
	if err != nil {
		log.Println("error create tx. shipment 18: " + err.Error())
		return err
	}
	result, err := db.GetAllOrdersForShipment(tx)
	if err != nil {
		log.Println("error get all orders for shipment: " + err.Error())
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Println("Emergency! Error in transaction!")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, result)
}

func GetPalletShipment(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	result, err := db.GetShipmentPalletModel(ctx.DB(), orderID)
	if err != nil {
		log.Println("Error get all palets for order: " + err.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func FinishPalletShipment(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	tx, err := ctx.DB().Begin()
	if err != nil {
		log.Println("error create tx. shipment 55: " + err.Error())
		return err
	}

	err = db.ShipTheOrder(tx, orderID)
	if err != nil {
		log.Println("error ship order: " + err.Error())
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Println("Emergency! Error in transaction!")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func GetPalletShipmentReport(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	orderID = orderID

	result := db.ShipmentReportModel{
		OrderCaption: "О-20-1325-РОССТАТ",
		Address:      "Какой-то там адрес",
		TotalBoxes:   206,
		TotalPallets: 2,
		Items: []db.ShipmentReportItemModel{
			{
				Num:                 1,
				Name:                "Форма №1. Бла балбалабала",
				Run:                 234,
				AmountInBox:         30,
				CompletedBoxes:      11,
				AmountInComposedBox: 14,
			},
			{
				Num:                 2,
				Name:                "Форма №2. Бла балбалабала",
				Run:                 234,
				AmountInBox:         50,
				CompletedBoxes:      4,
				AmountInComposedBox: 34,
			},
			{
				Num:                 3,
				Name:                "Форма №3. Бла балбалабала",
				Run:                 53,
				AmountInBox:         50,
				CompletedBoxes:      1,
				AmountInComposedBox: 3,
			},
		},
	}

	return ctx.JSON(http.StatusOK, result)
}
