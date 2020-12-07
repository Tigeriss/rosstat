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

	result, err := db.GetAllOrdersForShipment(ctx.DB())
	if err != nil {
		log.Println("error get all orders for shipment: " + err.Error())
		return err
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

	err = db.ShipTheOrder(ctx.DB(), orderID)
	if err != nil{
		log.Println("error ship order: " + err.Error())
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
