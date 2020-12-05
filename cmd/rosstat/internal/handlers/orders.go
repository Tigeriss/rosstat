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
	result, err := db.GetAllOrdersForCompletion(ctx.DB())
	if err != nil {
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

	// get the data by orderID
	result, err := db.GetOrderListForBigSuborder(ctx.DB(), orderID)
	if err != nil {
		log.Println("error GetOrderListForBigSuborder: " + err.Error())
		return err
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
	// get the data by orderID
	result, err := db.GetOrderListForSmallSuborder(ctx.DB(), orderID)
	if err != nil {
		log.Println("error GetOrderListForSmallSuborder: " + err.Error())
		return err
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

// GET /orders/big/pallet/:id - from /order/pallet page

func GetBigPalletOrders(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	// get the data by orderID
	result, err := db.GetOrderListForPallets(ctx.DB(), orderID)
	// 	db.BigPalletModel{
	// 	PalletNum: 3,
	// 	Types: []db.BigOrdersModel{
	// 		{
	// 			Type:     1,
	// 			FormName: "Форма №1. Записная книжечка переписчика (бла бла бла балб лабла бал)",
	// 		},
	// 		{
	// 			Type:     2,
	// 			FormName: "Форма №1. Записная книжечка Котофея Матвеевича",
	// 		},
	// 		{
	// 			Type:     3,
	// 			FormName: "Форма №1. Записная книжечка Выгебало",
	// 		},
	// 		{
	// 			Type:     4,
	// 			FormName: "Форма №1. Записная книжечка кадавра",
	// 		},
	// 		{
	// 			Type:     4,
	// 			FormName: "Форма №1. Записная книжечка кадавра",
	// 		},
	// 		{
	// 			Type:     4,
	// 			FormName: "Форма №1. Записная книжечка кадавра",
	// 		},
	// 	},
	// }

	return ctx.JSON(http.StatusOK, result)
}

func GetBigPalletNum(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	num, err := strconv.Atoi(c.Param("num"))
	if err != nil {
		return err
	}

	log.Println(num)
	log.Println(orderID)

	result := db.PrintPalletModel{
		OrderCaption:   "О-20-123-РОССТАТ 2",
		Address:        "107123, Москва",
		Provider:       "Жирпром",
		ContractNumber: "123-53322",
		Barcode:        "111222333",
		Register:       []db.PrintPalletRegisterModel{
			{
				NumPP:    1,
				Position: "Форма №2. Записная книжечка Котофея Матвеевича",
				Amount:   10,
				Boxes:    5,
			},
			{
				NumPP:    2,
				Position: "Форма №3. Записная книжечка Котофея Матвеевича",
				Amount:   10,
				Boxes:    5,
			},
			{
				NumPP:    3,
				Position: "Форма №4. Записная книжечка Котофея Матвеевича",
				Amount:   10,
				Boxes:    5,
			},
		},
	}

	return ctx.JSON(http.StatusOK, result)
}

// GET /orders/big/pallet/:id/barcode/:barcode

func GetBigPalletBarcodeOrders(c echo.Context) error {
	ctx := c.(*RosContext)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	barcode, err := strconv.Atoi(c.Param("barcode"))
	if err != nil {
		return err
	}

	orderID = orderID

	var result db.BigPalletBarcodeModel
	if barcode < 100 {
		result = db.BigPalletBarcodeModel{
			Success: false,
			Error:   "Товар с таким штрих-кодом не найден",
		}
	} else {
		result = db.BigPalletBarcodeModel{
			Success: true,
			Type:    barcode % 10,
		}
	}

	return ctx.JSON(http.StatusOK, result)
}

// POST /orders/big/pallet/:id/finish

func FinishBigPalletOrders(c echo.Context) error {
	ctx := c.(*RosContext)
	req := new(db.BigPalletFinishRequestModel)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := ctx.Bind(req); err != nil {
		return err
	}

	log.Println(orderID)
	log.Println(req)

	var result db.BigPalletFinishResponseModel
	if len(req.Barcodes) > 0 && req.Barcodes[0] == "111" {
		result = db.BigPalletFinishResponseModel{
			Success: false,
			Error:   "Невозможно завершить паллету",
		}
	} else if len(req.Barcodes) > 0 && req.Barcodes[0] == "221" {
		result = db.BigPalletFinishResponseModel{
			Success:    true,
			LastPallet: false,
		}
	} else {
		result = db.BigPalletFinishResponseModel{
			Success:    true,
			LastPallet: true,
		}
	}

	return ctx.JSON(http.StatusOK, result)
}
