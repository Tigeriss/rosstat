package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

var traceEnabled = true

func trace(name string, started time.Time) {
	if traceEnabled {
		log.Printf("Call: %s() - %s\n", name, time.Now().Sub(started))
	}
}

// for /orders
func GetAllOrdersForCompletion(tx *sql.Tx) ([]OrdersModel, error) {
	defer trace("GetAllOrdersForCompletion", time.Now())

	var result []OrdersModel

	statementGetOrders := "select id, num_order, contract, run, customer, order_name, address from rosstat.rosstat_orders where completed = false order by id;"
	rows, err := tx.Query(statementGetOrders)
	if err != nil {
		log.Println("error query statementGetOrders - select all orders")
		return nil, err
	}

	var index = 1
	for rows.Next() {
		order := Order{}
		err = rows.Scan(&order.Id, &order.NumOrder, &order.Contract, &order.Run, &order.Customer, &order.OrderName, &order.Address)
		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return nil, err
		}

		boxes, err := getCompletedBoxesAmountForOrder(tx, order.Id)
		if err != nil {
			log.Println("error get amount of box for order: " + err.Error())
		}

		pallets, err := getPalletsAmountForOrder(tx, order.Id)
		if err != nil {
			log.Println("error get amount of pallets for order: " + err.Error())
		}

		smallBoxes, err := getSmallBoxesAmountForOrder(tx, order.Id)
		if err != nil {
			log.Println("error get amount of combined boxes for order: " + err.Error())
		}

		tmp := OrdersModel{
			ID:            order.Id,
			Num:           index,
			OrderCaption:  order.NumOrder + "-" + order.OrderName,
			Customer:      order.Customer,
			Address:       order.Address,
			Run:           order.Run,
			AmountPallets: pallets,
			AmountBoxes:   boxes + smallBoxes,
			SubOrders: []SubOrderModel{
				{
					IsSmall:       false,
					OrderCaption:  order.NumOrder + "-" + order.OrderName + " короба",
					AmountPallets: pallets,
					AmountBoxes:   boxes,
				},
				{
					IsSmall:       true,
					OrderCaption:  order.NumOrder + "-" + order.OrderName + " сборные",
					AmountPallets: 0,
					AmountBoxes:   smallBoxes,
				},
			},
		}
		result = append(result, tmp)
		index++
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return nil, err
	}
	return result, nil
}

// for /shipment
func GetAllOrdersForShipment(tx *sql.Tx) ([]OrdersModel, error) {
	defer trace("GetAllOrdersForShipment", time.Now())

	var result []OrdersModel

	statementGetOrders := "select id, num_order, contract, run, customer, order_name, address from rosstat.rosstat_orders where completed = true and shipped = false  order by id;"
	rows, err := tx.Query(statementGetOrders)
	if err != nil {
		log.Println("error query statementGetOrders - select all orders")
		return nil, err
	}

	var index = 1
	for rows.Next() {
		order := Order{}
		err = rows.Scan(&order.Id, &order.NumOrder, &order.Contract, &order.Run, &order.Customer, &order.OrderName, &order.Address)
		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return nil, err
		}

		boxes, err := getCompletedBoxesAmountForOrder(tx, order.Id)
		if err != nil {
			log.Println("error get amount of box for order: " + err.Error())
		}

		pallets, err := getPalletsAmountForOrder(tx, order.Id)
		if err != nil {
			log.Println("error get amount of pallets for order: " + err.Error())
		}

		smallBoxes, err := getSmallBoxesAmountForOrder(tx, order.Id)
		if err != nil {
			log.Println("error get amount of combined boxes for order: " + err.Error())
		}

		tmp := OrdersModel{
			ID:            order.Id,
			Num:           index,
			OrderCaption:  order.NumOrder + "-" + order.OrderName,
			Customer:      order.Customer,
			Address:       order.Address,
			Run:           order.Run,
			AmountPallets: pallets,
			AmountBoxes:   boxes,
			SubOrders: []SubOrderModel{
				{
					IsSmall:       false,
					OrderCaption:  order.NumOrder + "-" + order.OrderName + " короба",
					AmountPallets: pallets,
					AmountBoxes:   boxes,
				},
				{
					IsSmall:       true,
					OrderCaption:  order.NumOrder + "-" + order.OrderName + " сборные",
					AmountPallets: 0,
					AmountBoxes:   smallBoxes,
				},
			},
		}
		result = append(result, tmp)
		index++
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return nil, err
	}
	return result, nil
}

func GetShipmentPalletModel(db *sql.DB, orderId int)([]ShipmentPalletModel, error){
	var result []ShipmentPalletModel
	allPalletInfos, err := GetAllPalletsBarcodesAndNums(db, orderId)
	if err != nil{
		log.Println("error get all pallets barcodes for order")
		return nil, err
	}
	for i := 0; i < len(allPalletInfos); i++{
		result = append(result, ShipmentPalletModel{
			Num:         i + 1,
			PalletNum:   allPalletInfos[i].palletNum,
			Barcode:     allPalletInfos[i].barcode,
			AmountBoxes: allPalletInfos[i].boxes,
		})
	}

	return result, nil
}

// for /orders/big
func GetOrderListForBigSuborder(tx *sql.Tx, orderId int) ([]BigOrdersModel, error) {
	var result []BigOrdersModel
	var amounts [27]int

	total, err := GetTotalBoxesAmount(tx, orderId)
	if err != nil {
		log.Println("error get boxes amount for order: " + err.Error())
		return nil, err
	}
	for i := 0; i < 26; i++ {
		amounts[i] = total[i]
	}

	// get amount of combined boxes to complete
	statement := "select count(id) from rosstat.small_boxes where order_id = " + strconv.Itoa(orderId) + ";"
	rows, err := tx.Query(statement)
	if err != nil {
		log.Println("error get amount of goods for order: " + err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&amounts[26])
		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return nil, err
		}
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return nil, err
	}

	var allCompletedBoxes [27]int
	allCompletedBoxIds, err := GetCompletedBoxesIds(tx, orderId)
	if len(allCompletedBoxIds) != 0 {
		for i := 0; i < len(allCompletedBoxIds); i++ {
			good := GetProductByBoxID(allCompletedBoxIds[i])
			allCompletedBoxes[good.Type] ++
		}
	}

	for i := 1; i < 27; i++ {
		if amounts[i-1] != 0 {
			good := GetProductByType(i)
			boxes := allCompletedBoxes[i-1]
			if err != nil {
				log.Println("error get completed boxes amount of certain product for order: " + err.Error())
				return nil, err
			}

			tmp := BigOrdersModel{}
			tmp.FormName = good.Name
			tmp.Type = good.Type
			tmp.Total = total[i-1]
			tmp.Built = boxes
			result = append(result, tmp)
		}

	}
	return result, nil
}

// for /orders/small
func GetOrderListForSmallSuborder(db *sql.DB, orderId int) ([]BigOrdersModel, error) {
	// yep, its copy-paste, but let it be for now
	var result []BigOrdersModel
	var amounts [26]int
	statement := "select "
	for i := 1; i < 27; i++ {
		statement += "\"" + strconv.Itoa(i) + "\","
	}
	statement = strings.TrimRight(statement, ",")
	statement += " from rosstat.rosstat_orders where id = " + strconv.Itoa(orderId) + ";"

	// get total amount of every good type
	rows, err := db.Query(statement)
	if err != nil {
		log.Println("error get amount of goods for order: " + err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&amounts[0], &amounts[1], &amounts[2],
			&amounts[3], &amounts[4], &amounts[5],
			&amounts[6], &amounts[7], &amounts[8],
			&amounts[9], &amounts[10], &amounts[11],
			&amounts[12], &amounts[13], &amounts[14],
			&amounts[15], &amounts[16], &amounts[17],
			&amounts[18], &amounts[19], &amounts[20],
			&amounts[21], &amounts[22], &amounts[23],
			&amounts[24], &amounts[25])
		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return nil, err
		}
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return result, err
	}
	allProductsTotal, err := GetTotalPiecesAmountForOrder(db, orderId)
	if err != nil {
		log.Println("error get total pieces for order: " + err.Error())
		return result, err
	}

	for i := 1; i < 27; i++ {
		if amounts[i-1] != 0 {

			good := GetProductByType(i)
			if err != nil {
				log.Println("error get boxes amount of certain product for order: " + err.Error())
				return nil, err
			}

			tmp := BigOrdersModel{}
			tmp.FormName = good.Name
			tmp.Type = good.Type
			tmp.Total = allProductsTotal[i-1]
			tmp.Built = 0
			result = append(result, tmp)
		}

	}
	return result, nil
}

// for /orders/pallet
func GetOrderListForPallets(tx *sql.Tx, orderId int) (BigPalletModel, error) {
	var result BigPalletModel
	var allPalletsNums []int
	palNum := 0
	// get num for pallet:
	statementGetNum := "select num from rosstat.pallets where order_id = " +
		strconv.Itoa(orderId) + " order by num desc;"

	rows, err := tx.Query(statementGetNum)
	if err != nil {
		log.Println("error query statementGetNum - select last pallet num")
		return BigPalletModel{}, err
	}

	if rows.Next() {
		err = rows.Scan(&palNum)
		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return BigPalletModel{}, err
		}
		allPalletsNums = append(allPalletsNums, palNum)
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return BigPalletModel{}, err
	}
	if len(allPalletsNums) > 0 {
		palNum = allPalletsNums[0]
	} else {
		palNum = 0
	}

	result.PalletNum = palNum + 1
	allBoxes, err := GetBoxesToCompleteForOrder(tx, orderId)
	for i := 0; i < len(allBoxes); i++ {
		good := GetProductByType(i + 1)
		for s := 0; s < allBoxes[i]; s++ {
			result.Types = append(result.Types, BigOrdersModel{
				Type:     good.Type,
				FormName: good.Name,
				Total:    0,
				Built:    0,
			})
		}
	}
	return result, nil
}

// Call it after button "combined boxes fully completed" pressed
func PutSmallOrderToDB(tx *sql.Tx, orderId int, boxIds []string, us string) (int, error) {

	err := createSmallBoxesRecord(tx, orderId, boxIds, us)
	if err != nil {
		log.Println("Error create record in rosstat.small_boxes: " + err.Error())
		return 0, err
	}

	return len(boxIds), nil
}

// for print pallet label and registry
func GetDataForLabelAndRegistry(db *sql.DB, orderId, palletNum int) (PrintPalletModel, error) {
	var result PrintPalletModel

	ord, err := GetOrderById(db, orderId)
	if err != nil {
		log.Println("error get order: " + err.Error())
		return result, err
	}

	palletId := fmt.Sprintf("%d%d", orderId, palletNum)
	result.OrderCaption = ord.NumOrder + "-" + ord.OrderName
	result.Address = ord.Address
	// TODO: ask about provider
	result.Provider = "ООО ББС"
	result.ContractNumber = ord.Contract

	result.Barcode = fmt.Sprintf("%012s", palletId)
	boxIdsForPallet, err := GetBoxIdsForPallet(db, palletId)
	if err != nil {
		log.Println("error get boxes for pallet: " + err.Error())
		return result, err
	}

	dataForRegistry := [27]PalletRegistryGoodsData{}
	for i := 0; i < len(boxIdsForPallet); i++{
		good := GetProductByBoxID(boxIdsForPallet[i])
		dataForRegistry[good.Type-1].good = good
		dataForRegistry[good.Type -1].boxes++
	}

	var reg []PrintPalletRegisterModel
	index := 1
	for i := 0; i < len(dataForRegistry); i++{
		var tmp = PrintPalletRegisterModel{}
		if dataForRegistry[i].boxes > 0{
			tmp.NumPP = index
			tmp.Position = dataForRegistry[i].good.Name
			tmp.Boxes = dataForRegistry[i].boxes
			tmp.Amount = dataForRegistry[i].boxes * dataForRegistry[i].good.AmountInBox
			reg = append(reg,tmp)
			index++
		}
	}
	result.Register = reg
	return result, nil
}

func GetDataForGetBigPalletBarcodeOrders(barcode int)BigPalletBarcodeModel{
	good := GetProductByBoxID(barcode)
	if good.Type != 0{
		return BigPalletBarcodeModel{
			Success: true,
			Type:    good.Type,
			Error:   "",
		}
	}else{
		return BigPalletBarcodeModel{
			Success: false,
			Type:    0,
			Error:   "Товара с таким штрихкодом не существует! Проверьте короб",
		}
	}
}

func GetOrderById(db *sql.DB, orderId int) (Order, error) {
	var result Order
	statementGetOrders := "select id, num_order, contract, run, customer, order_name, address from rosstat.rosstat_orders where id = " + strconv.Itoa(orderId) + ";"
	rows, err := db.Query(statementGetOrders)
	if err != nil {
		log.Println("error query statementGetOrders - select all orders")
		return result, err
	}

	for rows.Next() {
		err = rows.Scan(&result.Id, &result.NumOrder, &result.Contract, &result.Run, &result.Customer, &result.OrderName, &result.Address)
		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return result, err
		}
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return result, err
	}
	return result, nil
}

func GetBoxIdsForPallet(db *sql.DB, palletId string) ([]int, error){
	var result []int
	statement := "select id from rosstat.boxes where pallet_id = " + palletId + ";"
	rows, err := db.Query(statement)
	if err != nil{
		log.Println("error get ids for pallet: " + err.Error())
		return nil, err
	}
	id := 0
	for rows.Next(){
		rows.Scan(&id)
		result = append(result, id)
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return nil, err
	}
	return result, nil
}

func CreatePallet(tx *sql.Tx, orderId, palletNum int, boxes []string, us string) (BigPalletFinishResponseModel, error) {

	var result BigPalletFinishResponseModel
	palletId := strconv.Itoa(orderId) + strconv.Itoa(palletNum)
	isLast, isSmallCompleted := isPalletLastInOrder(tx, orderId, boxes)
	// 1. check if it's last pallet
	if isLast {
		// 2a. if true, check if small suborder in barcodes
		if isSmallCompleted {
			// 3a.a. if true - put data in db, update order as completed, return success, last - true
			err := createPalletRecord(tx, orderId, palletNum, palletId)
			if err != nil {
				log.Println("error create pallet record: " + err.Error())
				return BigPalletFinishResponseModel{
					Success:    false,
					Error:      "Не удалось создать запись паллеты в базе данных",
					LastPallet: true,
				}, err
			}
			err = createBoxesRecord(tx, palletId, boxes, us)
			if err != nil {
				log.Println("error create boxes record: " + err.Error())
				return BigPalletFinishResponseModel{
					Success:    false,
					Error:      "Не удалось создать запись коробов в базе данных",
					LastPallet: true,
				}, err
			}
			err = completeTheOrder(tx, orderId)
			if err != nil {
				log.Println("error update order to completed: " + err.Error())
				return BigPalletFinishResponseModel{
					Success:    false,
					Error:      "Не удалось обновить статус заказа в базе данных",
					LastPallet: true,
				}, err
			}
			result.Error = ""
			result.LastPallet = isLast
			result.Success = true

		} else {
			// 3a.b. if false return error
			result.Success = false
			result.LastPallet = true
			result.Error = "Внимание. Это последняя паллета, но у вас не собран малый подзаказ. Паллета не может быть завершена!"
		}

	} else {
		// 2b. if false - put data in db, return success, last - false
		err := createPalletRecord(tx, orderId, palletNum, palletId)
		if err != nil {
			log.Println("error create pallet record: " + err.Error())
			return BigPalletFinishResponseModel{
				Success:    false,
				Error:      "Не удалось создать запись паллеты в базе данных",
				LastPallet: false,
			}, err
		}
		err = createBoxesRecord(tx, palletId, boxes, us)
		if err != nil {
			log.Println("error create boxes record: " + err.Error())
			return BigPalletFinishResponseModel{
				Success:    false,
				Error:      "Не удалось создать запись коробов в базе данных",
				LastPallet: false,
			}, err
		}
		result.Error = ""
		result.LastPallet = false
		result.Success = true
	}

	return result, nil
}

func GetBoxesToCompleteForOrder(tx *sql.Tx, orderId int) ([]int, error) {
	var result []int
	var totalAmounts []int
	var completedBoxes []int
	combinedBoxes := 0

	totalAmounts, err := GetTotalBoxesAmount(tx, orderId)
	if err != nil {
		log.Println("error get total boxes amount: " + err.Error())
		return nil, err
	}
	completedBoxes, err = GetCompletedBoxesAmount(tx, orderId)
	if err != nil {
		log.Println("error get completed boxes amount: " + err.Error())
		return nil, err
	}
	combinedBoxes, err = getSmallBoxesAmountForOrder(tx, orderId)
	if err != nil {
		log.Println("error get combined boxes amount: " + err.Error())
		return nil, err
	}
	for i := 0; i < 26; i++ {
		result = append(result, totalAmounts[i]-completedBoxes[i])
	}

	result = append(result, combinedBoxes)

	return result, nil
}

func GetTotalPiecesAmountForOrder(db *sql.DB, orderId int) ([]int, error) {
	var amounts [26]int
	var result []int

	statement := "select "
	for i := 1; i < 27; i++ {
		statement += "\"" + strconv.Itoa(i) + "\","
	}
	statement = strings.TrimRight(statement, ",")
	statement += " from rosstat.rosstat_orders where id = " + strconv.Itoa(orderId) + ";"

	// get total amount of every good type
	rows, err := db.Query(statement)
	if err != nil {
		log.Println("error get total amount of goods for order: " + err.Error())
		return []int{0}, err
	}
	for rows.Next() {
		err = rows.Scan(&amounts[0], &amounts[1], &amounts[2],
			&amounts[3], &amounts[4], &amounts[5],
			&amounts[6], &amounts[7], &amounts[8],
			&amounts[9], &amounts[10], &amounts[11],
			&amounts[12], &amounts[13], &amounts[14],
			&amounts[15], &amounts[16], &amounts[17],
			&amounts[18], &amounts[19], &amounts[20],
			&amounts[21], &amounts[22], &amounts[23],
			&amounts[24], &amounts[25])

		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return []int{0}, err
		}
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return result, err
	}
	for i := 1; i < 27; i++ {
		result = append(result, GetPiecesOfThisProduct(GetProductByType(i), amounts[i-1]))
	}

	return result, nil
}

func GetTotalBoxesAmount(tx *sql.Tx, orderId int) ([]int, error) {
	var amounts [26]int
	var result []int

	statement := "select "
	for i := 1; i < 27; i++ {
		statement += "\"" + strconv.Itoa(i) + "\","
	}
	statement = strings.TrimRight(statement, ",")
	statement += " from rosstat.rosstat_orders where id = " + strconv.Itoa(orderId) + ";"
	// log.Println(statement)

	// get total amount of every good type
	rows, err := tx.Query(statement)
	if err != nil {
		log.Println("error get total amount of goods for order: " + err.Error())
		return []int{0}, err
	}
	for rows.Next() {
		err = rows.Scan(&amounts[0], &amounts[1], &amounts[2],
			&amounts[3], &amounts[4], &amounts[5],
			&amounts[6], &amounts[7], &amounts[8],
			&amounts[9], &amounts[10], &amounts[11],
			&amounts[12], &amounts[13], &amounts[14],
			&amounts[15], &amounts[16], &amounts[17],
			&amounts[18], &amounts[19], &amounts[20],
			&amounts[21], &amounts[22], &amounts[23],
			&amounts[24], &amounts[25])

		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return []int{0}, err
		}
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return result, err
	}
	for i := 1; i < 27; i++ {
		result = append(result, GetWholeBoxesOfThisProduct(GetProductByType(i), amounts[i-1]))
	}

	return result, nil
}

func GetCompletedBoxesAmount(tx *sql.Tx, orderId int) ([]int, error) {
	var result []int
	var amounts [26]int
	statement := "select id from rosstat.boxes where pallet_id in " +
		"(select id from rosstat.pallets where order_id = " + strconv.Itoa(orderId) + ") order by id;"
	rows, err := tx.Query(statement)
	if err != nil {
		log.Println("error query select all boxes amount")
		return nil, err
	}
	boxId := 0
	for rows.Next() {
		err = rows.Scan(&boxId)
		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return nil, err
		}
		good := GetProductByBoxID(boxId)
		amounts[good.Type-1] ++

	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return nil, err
	}
	for i := 0; i < 26; i++ {
		result = append(result, amounts[i])
	}
	return result, nil
}

func GetCompletedBoxesIds(tx *sql.Tx, orderId int) ([]int, error) {
	var result []int
	statement := "select id from rosstat.boxes where pallet_id in " +
		"(select id from rosstat.pallets where order_id = " + strconv.Itoa(orderId) + ");"

	rows, err := tx.Query(statement)
	if err != nil {
		log.Println("error query select all boxes ids")
		return nil, err
	}
	boxId := 0
	for rows.Next() {
		err = rows.Scan(&boxId)
		if err != nil {
			log.Println("error get data from row: " + err.Error())
			return nil, err
		}
		result = append(result, boxId)
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return nil, err
	}
	return result, nil
}

func getCompletedBoxesAmountForOrder(tx *sql.Tx, orderId int) (int, error) {
	// defer trace("getCompletedBoxesAmountForOrder", time.Now())

	var boxes = 0
	// var pallets []int
	// statement := "select id from rosstat.pallets where order_id = " + strconv.Itoa(orderId) + ";"
	// rows, err := tx.Query(statement)
	// if err != nil {
	// 	log.Println("error get amount of pallets for order: " + err.Error())
	// }
	// pall := 0
	// for rows.Next() {
	// 	rows.Scan(&pall)
	// 	pallets = append(pallets, pall)
	// }
	// rows.Close()
	// if len(pallets) != 0{
	// 	statement := "select count(id) from rosstat.boxes where pallet_id in ("
	// 	for i := 0; i < len(pallets); i++{
	//
	// 	}
	// }

	statement := "select count(id) from rosstat.boxes where pallet_id in " +
		"(select id from rosstat.pallets where order_id = " + strconv.Itoa(orderId) + ");"
	rows, err := tx.Query(statement)
	if err != nil {
		log.Println("error get amount of boxes for order: " + err.Error())
	}
	for rows.Next() {
		rows.Scan(&boxes)
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return 0, err
	}

	return boxes, nil
}

func getPalletsAmountForOrder(tx *sql.Tx, orderId int) (int, error) {
	// defer trace("getPalletsAmountForOrder", time.Now())

	var pallets = 0

	statement := "select count(id) from rosstat.pallets where order_id = " + strconv.Itoa(orderId) + ";"
	rows, err := tx.Query(statement)
	if err != nil {
		log.Println("error get amount of pallets for order: " + err.Error())
	}
	for rows.Next() {
		rows.Scan(&pallets)
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return 0, err
	}
	return pallets, nil
}

func getSmallBoxesAmountForOrder(tx *sql.Tx, orderId int) (int, error) {
	// defer trace("getSmallBoxesAmountForOrder", time.Now())

	var boxes = 0

	statement := "select count(id) from rosstat.small_boxes where order_id = " + strconv.Itoa(orderId) + ";"
	rows, err := tx.Query(statement)
	if err != nil {
		log.Println("error get amount of small boxes for order: " + err.Error())
	}
	for rows.Next() {
		rows.Scan(&boxes)
	}
	err = rows.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return boxes, err
	}
	return boxes, nil
}

// get amount of whole boxes of particular good
func GetWholeBoxesOfThisProduct(product Good, totalAmount int) int {
	return totalAmount / product.AmountInBox
}

// get amount of pieces of particular good for combined box
func GetPiecesOfThisProduct(product Good, totalAmount int) int {
	// log.Println("total: " + strconv.Itoa(totalAmount) + ", Amount in a box: " + strconv.Itoa(product.AmountInBox))
	return totalAmount % product.AmountInBox
}

func GetProductByBoxID(id int) Good {
	result := Good{
	}

	if id >= 200100001 && id <= 200200000 {
		result.Name = "1. Форма № 1. Записная книжка переписчика (является приложением к Инструкции)"
		result.Run = 476596
		result.AmountInBox = 20
		result.Type = 1

	} else if id >= 200200001 && id <= 200300000 {
		result.Name = "2. Форма № 2. Записная книжка контролера полевого уровня"
		result.Run = 65357
		result.AmountInBox = 50
		result.Type = 2

	} else if id >= 200300001 && id <= 200400000 {
		result.Name = "3. Форма № 3. Записная книжка уполномоченного по вопросам переписи"
		result.Run = 6023
		result.AmountInBox = 50
		result.Type = 3

	} else if id >= 200400001 && id <= 200500000 {
		result.Name = "4. Форма № 4. Сводная ведомость по переписному участку"
		result.Run = 57930
		result.AmountInBox = 1000
		result.Type = 4

	} else if id >= 200500001 && id <= 200600000 {
		result.Name = "5. Форма № 5. Сводная ведомость по городскому округу, муниципальному району/ округу"
		result.Run = 10459
		result.AmountInBox = 1000
		result.Type = 5

	} else if id >= 200600001 && id <= 200700000 {
		result.Name = "6. Форма № 6. Сводка итогов переписи населения по городскому округу, муниципальному району/округу"
		result.Run = 61459
		result.AmountInBox = 500
		result.Type = 6

	} else if id >= 200700001 && id <= 200800000 {
		result.Name = "7. Форма № 7. Информационные листовки (к лицам, которых трудно застать дома)"
		result.Run = 28419540
		result.AmountInBox = 2000
		result.Type = 7

	} else if id >= 200800001 && id <= 200900000 {
		result.Name = "8. Форма № 9. Ярлык в портфель переписчика"
		result.Run = 18812
		result.AmountInBox = 8000
		result.Type = 8

	} else if id >= 200900001 && id <= 201000000 {
		result.Name = "9. Форма № 10. Карточка для респондентов"
		result.Run = 392150
		result.AmountInBox = 2000
		result.Type = 9

	} else if id >= 201000001 && id <= 201100000 {
		result.Name = "10. Форма Обложка. Обложка на переписные документы"
		result.Run = 2287448
		result.AmountInBox = 500
		result.Type = 10

	} else if id >= 201100001 && id <= 201200000 {
		result.Name = "11. Форма С. Список лиц"
		result.Run = 2287448
		result.AmountInBox = 1000
		result.Type = 11

	} else if id >= 201200001 && id <= 201300000 {
		result.Name = "12. Форма КС. Список лиц для контроля за заполнением переписных листов"
		result.Run = 790907
		result.AmountInBox = 2000
		result.Type = 12

	} else if id >= 201300001 && id <= 201400000 {
		result.Name = "13. Форма СПР. Справка о прохождении переписи"
		result.Run = 10136033
		result.AmountInBox = 8000
		result.Type = 13

	} else if id >= 201400001 && id <= 201500000 {
		result.Name = "14. Инструкция о порядке подготовки материалов Всероссийской переписи населения 2020 года к обработке"
		result.Run = 1652
		result.AmountInBox = 40
		result.Type = 14

	} else if id >= 201500001 && id <= 201600000 {
		result.Name = "15. Тесты для обучения переписного персонала"
		result.Run = 495982
		result.AmountInBox = 100
		result.Type = 15

	} else if id >= 201600001 && id <= 201700000 {
		result.Name = "16. Указатели для переписных участков"
		result.Run = 32689
		result.AmountInBox = 500
		result.Type = 16

	} else if id >= 201700001 && id <= 201800000 {
		result.Name = "17. Форма Л. Переписной лист (обучение, чистые формы)"
		result.Run = 2082395
		result.AmountInBox = 1000
		result.Type = 17

	} else if id >= 201800001 && id <= 201900000 {
		result.Name = "18. Форма П. Переписной лист (обучение, чистые формы)"
		result.Run = 832958
		result.AmountInBox = 1000
		result.Type = 18

	} else if id >= 201900001 && id <= 202000000 {
		result.Name = "19. Форма В. Переписной лист (обучение, чистые формы)"
		result.Run = 416479
		result.AmountInBox = 1000
		result.Type = 19

	} else if id >= 202000001 && id <= 202100000 {
		result.Name = "20. Форма Н. Сопроводительный бланк (обучение, чистые формы)"
		result.Run = 416479
		result.AmountInBox = 1000
		result.Type = 20

	} else if id >= 202100001 && id <= 202200000 {
		result.Name = "21. Форма Обложка. Обложка на переписные документы (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 500
		result.Type = 21

	} else if id >= 202200001 && id <= 202300000 {
		result.Name = "22. Форма С. Список лиц (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Type = 22

	} else if id >= 202300001 && id <= 202400000 {
		result.Name = "23. Форма Л. Переписной лист (обучение, заполненные формы)"
		result.Run = 8260
		result.AmountInBox = 1000
		result.Type = 23

	} else if id >= 202400001 && id <= 202500000 {
		result.Name = "24. Форма П. Переписной лист (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Type = 24

	} else if id >= 202500001 && id <= 202600000 {
		result.Name = "25. Форма В. Переписной лист (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Type = 25

	} else if id >= 202600001 && id <= 202700000 {
		result.Name = "26. Форма Н. Сопроводительный бланк (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Type = 26

	} else if id >= 202700001 && id <= 202800000 {
		result.Name = "27. Сборный короб "
		result.Run = 0
		result.AmountInBox = 0
		result.Type = 27

	} else{
		result.Name = "Ошибка!"
		result.Run = 0
		result.AmountInBox = 0
		result.Type = 0
	}

	return result
}

func GetProductByType(t int) Good {
	result := Good{
	}

	switch t {
	case 1:
		result.Name = "1. Форма № 1. Записная книжка переписчика (является приложением к Инструкции)"
		result.Run = 476596
		result.AmountInBox = 20
		result.Type = 1
		result.FirstID = 200100001
		result.LastID = 200200000
		break
	case 2:
		result.Name = "2. Форма № 2. Записная книжка контролера полевого уровня"
		result.Run = 65357
		result.AmountInBox = 50
		result.Type = 2
		result.FirstID = 200200001
		result.LastID = 200300000
		break
	case 3:
		result.Name = "3. Форма № 3. Записная книжка уполномоченного по вопросам переписи"
		result.Run = 6023
		result.AmountInBox = 50
		result.Type = 3
		result.FirstID = 200300001
		result.LastID = 200400000
		break
	case 4:
		result.Name = "4. Форма № 4. Сводная ведомость по переписному участку"
		result.Run = 57930
		result.AmountInBox = 1000
		result.Type = 4
		result.FirstID = 200400001
		result.LastID = 200500000
		break
	case 5:
		result.Name = "5. Форма № 5. Сводная ведомость по городскому округу, муниципальному району/ округу"
		result.Run = 10459
		result.AmountInBox = 1000
		result.Type = 5
		result.FirstID = 200500001
		result.LastID = 200600000
		break
	case 6:
		result.Name = "6. Форма № 6. Сводка итогов переписи населения по городскому округу, муниципальному району/округу"
		result.Run = 61459
		result.AmountInBox = 500
		result.Type = 6
		result.FirstID = 200600001
		result.LastID = 200700000
		break
	case 7:
		result.Name = "7. Форма № 7. Информационные листовки (к лицам, которых трудно застать дома)"
		result.Run = 28419540
		result.AmountInBox = 2000
		result.Type = 7
		result.FirstID = 200700001
		result.LastID = 200800000
		break
	case 8:
		result.Name = "8. Форма № 9. Ярлык в портфель переписчика"
		result.Run = 18812
		result.AmountInBox = 8000
		result.Type = 8
		result.FirstID = 200800001
		result.LastID = 200900000
		break
	case 9:
		result.Name = "9. Форма № 10. Карточка для респондентов"
		result.Run = 392150
		result.AmountInBox = 2000
		result.Type = 9
		result.FirstID = 200900001
		result.LastID = 201000000
		break
	case 10:
		result.Name = "10. Форма Обложка. Обложка на переписные документы"
		result.Run = 2287448
		result.AmountInBox = 500
		result.Type = 10
		result.FirstID = 201000001
		result.LastID = 201100000
		break
	case 11:
		result.Name = "11. Форма С. Список лиц"
		result.Run = 2287448
		result.AmountInBox = 1000
		result.Type = 11
		result.FirstID = 201100001
		result.LastID = 201200000
		break
	case 12:
		result.Name = "12. Форма КС. Список лиц для контроля за заполнением переписных листов"
		result.Run = 790907
		result.AmountInBox = 2000
		result.Type = 12
		result.FirstID = 201200001
		result.LastID = 201300000
		break
	case 13:
		result.Name = "13. Форма СПР. Справка о прохождении переписи"
		result.Run = 10136033
		result.AmountInBox = 8000
		result.Type = 13
		result.FirstID = 201300001
		result.LastID = 201400000
		break
	case 14:
		result.Name = " 14. Инструкция о порядке подготовки материалов Всероссийской переписи населения 2020 года к обработке"
		result.Run = 1652
		result.AmountInBox = 40
		result.Type = 14
		result.FirstID = 201400001
		result.LastID = 201500000
		break
	case 15:
		result.Name = "15. Тесты для обучения переписного персонала"
		result.Run = 495982
		result.AmountInBox = 100
		result.Type = 15
		result.FirstID = 201500001
		result.LastID = 201600000
		break
	case 16:
		result.Name = "16. Указатели для переписных участков"
		result.Run = 32689
		result.AmountInBox = 500
		result.Type = 16
		result.FirstID = 201600001
		result.LastID = 201700000
		break
	case 17:
		result.Name = "17. Форма Л. Переписной лист (обучение, чистые формы)"
		result.Run = 2082395
		result.AmountInBox = 1000
		result.Type = 17
		result.FirstID = 201700001
		result.LastID = 201800000
		break
	case 18:
		result.Name = "18. Форма П. Переписной лист (обучение, чистые формы)"
		result.Run = 832958
		result.AmountInBox = 1000
		result.Type = 18
		result.FirstID = 201800001
		result.LastID = 201900000
		break
	case 19:
		result.Name = "19. Форма В. Переписной лист (обучение, чистые формы)"
		result.Run = 416479
		result.AmountInBox = 1000
		result.Type = 19
		result.FirstID = 201900001
		result.LastID = 202000000
		break
	case 20:
		result.Name = "20. Форма Н. Сопроводительный бланк (обучение, чистые формы)"
		result.Run = 416479
		result.AmountInBox = 1000
		result.Type = 20
		result.FirstID = 202000001
		result.LastID = 202100000
		break
	case 21:
		result.Name = "21. Форма Обложка. Обложка на переписные документы (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 500
		result.Type = 21
		result.FirstID = 202100001
		result.LastID = 202200000
		break
	case 22:
		result.Name = "22. Форма С. Список лиц (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Type = 22
		result.FirstID = 202200001
		result.LastID = 202300000
		break
	case 23:
		result.Name = "23. Форма Л. Переписной лист (обучение, заполненные формы)"
		result.Run = 8260
		result.AmountInBox = 1000
		result.Type = 23
		result.FirstID = 202300001
		result.LastID = 202400000
		break
	case 24:
		result.Name = "24. Форма П. Переписной лист (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Type = 24
		result.FirstID = 202400001
		result.LastID = 202500000
		break
	case 25:
		result.Name = "25. Форма В. Переписной лист (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Type = 25
		result.FirstID = 202500001
		result.LastID = 202600000
		break
	case 26:
		result.Name = "26. Форма Н. Сопроводительный бланк (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Type = 26
		result.FirstID = 202600001
		result.LastID = 202700000
		break
	case 27:
		result.Name = "27. Сборный короб"
		result.Run = 0
		result.AmountInBox = 0
		result.Type = 27
		result.FirstID = 202700001
		result.LastID = 202800000
		break
	default:
		result.Name = "Ошибка!"
		result.Run = 0
		result.AmountInBox = 0
		result.Type = 0
		break

	}

	return result
}

func GetAllPalletsBarcodesAndNums(db *sql.DB, orderId int)([]PalletInfo, error){
	var result []PalletInfo
	statement1 := "select id,num from rosstat.pallets where order_id = " + strconv.Itoa(orderId) + " order by id;"
	rows1, err := db.Query(statement1)
	if err != nil {
		log.Println("error get all pallet ids for order: " + err.Error())
		return nil, err
	}
	statement2 := "select count(id), pallet_id from rosstat.boxes where pallet_id in " +
		"(select id from rosstat.pallets where order_id =" + strconv.Itoa(orderId) +
		" ) group by rosstat.boxes.pallet_id order by rosstat.boxes.pallet_id;"
	rows2, err := db.Query(statement2)
	if err != nil {
		log.Println("error get all boxes for order: " + err.Error())
		return nil, err
	}
	for rows1.Next(){
		tmpId := ""
		tmpNum := 0
		tmpBoxes := 0
		err = rows1.Scan(&tmpId, &tmpNum)
		if err != nil {
			log.Println("error scan pallet id for order: " + err.Error())
			return nil, err
		}
		bar, err := generateBarcodeWithControlNumForPalletId(tmpId)
		if err != nil {
			log.Println("error generate barcode for pallet id: " + err.Error())
			return nil, err
		}
		result = append(result, PalletInfo{
			barcode:  bar,
			palletNum: tmpNum,
			boxes: tmpBoxes,
		})
	}
	ind := 0
	for rows2.Next(){
		tmpBoxes := 0
		tmpPalletId := 0
		rows2.Scan(&tmpBoxes, &tmpPalletId)
		result[ind].boxes = tmpBoxes
		ind ++
	}
	err = rows1.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return nil, err
	}
	err = rows2.Close()
	if err != nil {
		log.Println("error close row: " + err.Error())
		return nil, err
	}
	return result, nil
}

// put data in rosstat.boxes
func createBoxesRecord(tx *sql.Tx, palletId string, boxes []string, us string) error {
	statement := "insert into rosstat.boxes values " // rosstat.boxes columns: id, pallet_id, us_name
	for i := 0; i < len(boxes); i++ {
		statement += "(" + boxes[i] + ", " + palletId + ", '" + us + "'),"
	}
	statement = strings.TrimRight(statement, ",")
	statement += ";"
	_, err := tx.Exec(statement)
	if err != nil {
		log.Println("error execute query to insert boxes for order")
		if err := tx.Rollback(); err != nil {
			log.Println("We were unable to rollback transaction. That's odd but we really can't do anything else here")
		}
		return err
	}
	return nil
}

// put data in rosstat.small_boxes
func createSmallBoxesRecord(tx *sql.Tx, orderId int, boxIds []string, us string) error {
	log.Println(us + " in createSmallBoxesRecord")
	statementInsert := "insert into rosstat.small_boxes values "
	for i := 0; i < len(boxIds); i++ {
		statementInsert += "(" + boxIds[i] + ", " + strconv.Itoa(orderId) + ", '" + us + "'),"
	}
	statementInsert = strings.TrimRight(statementInsert, ",")
	statementInsert += ";"
	log.Println("st: " + statementInsert)
	_, err := tx.Exec(statementInsert)
	if err != nil {
		log.Println("error execute query to insert small boxes for order. Transaction will be roll back")
		if err := tx.Rollback(); err != nil {
			log.Println("We were unable to rollback transaction. That's odd but we really can't do anything else here")
		}
		return err
	}
	return nil
}

// check if pallet last in order returns: is last, is small completed
func isPalletLastInOrder(tx *sql.Tx, orderId int, barcodes []string) (bool, bool) {
	isLast := false
	isSmallCompleted := false
	// 1. get total boxes for order
	total, err := GetTotalBoxesAmount(tx, orderId)
	if err != nil {
		log.Println("error get total amount of boxes for order")
	}
	tot := 0
	for i := 0; i < len(total); i++ {
		tot += total[i]
	}
	// 2. get completed boxes for order
	completed, err := GetCompletedBoxesAmount(tx, orderId)
	if err != nil {
		log.Println("error get amount of completed boxes for order")
	}
	comp := 0
	for i := 0; i < len(completed); i++ {
		comp += completed[i]
	}
	// 3. get small boxes for order
	small, err := getSmallBoxesAmountForOrder(tx, orderId)
	if err != nil {
		log.Println("error get amount of small boxes for order")
	}
	// 4. check if ( (total + small) - completed) > len(barcodes)
	if ((tot + small) - comp) <= len(barcodes) {
		isLast = true
	}
	// 5. check if small boxes exists
	if small > 0 {
		isSmallCompleted = true
	}

	return isLast, isSmallCompleted
}

// put pallet data to db
func createPalletRecord(tx *sql.Tx, orderId, palletNum int, palletId string) error {
	statement := "insert into rosstat.pallets values(" +
		palletId + ", " +
		strconv.Itoa(palletNum) + ", " + strconv.Itoa(orderId) + ");"
	_, err := tx.Exec(statement)
	if err != nil {
		log.Println("error insert into pallets: " + err.Error())
		if err := tx.Rollback(); err != nil {
			log.Println("Transaction rollback failed!")
		}
		return err
	}
	return nil
}

func completeTheOrder(tx *sql.Tx, orderId int) error {

	statement := "update rosstat.rosstat_orders set completed = true where id = " + strconv.Itoa(orderId) + ";"
	_, err := tx.Exec(statement)
	if err != nil {
		log.Println("error update completion of order: " + err.Error())
		if err := tx.Rollback(); err != nil {
			log.Println("Transaction rollback failed!")
		}
		return err
	}

	return nil
}

func ShipTheOrder(tx *sql.Tx, orderId int) error {

	statement := "update rosstat.rosstat_orders set shipped = true where id = " + strconv.Itoa(orderId) + ";"
	_, err := tx.Exec(statement)
	if err != nil {
		log.Println("error update completion of order: " + err.Error())
		if err := tx.Rollback(); err != nil {
			log.Println("Transaction rollback failed!")
		}
		return err
	}

	return nil
}

func generateBarcodeWithControlNumForPalletId(palletId string) (string, error){
	palletId = fmt.Sprintf("%012s", palletId)
	odd := 0
	even := 0
	for i:= 0; i < len(palletId); i++{
		n, err := strconv.Atoi(string(palletId[i]))
		if err != nil {
			log.Println("error convert string to int: " + err.Error())
			return "", err
		}
		if i == 0 || i %2 != 0{
			even += n
		}else{
			odd += n
		}
	}
	even *= 3
	res := even + odd
	res = res % 10
	conNum := 10 - res

	return palletId + strconv.Itoa(conNum), nil
}