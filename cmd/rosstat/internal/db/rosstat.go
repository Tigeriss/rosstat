package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
)

// good itself. hardcoded because of problems with id
type Good struct {
	Num         string `db:"num" json:"num"`
	Name        string `db:"name" json:"name"`
	Run         int    `db:"run" json:"run"`
	AmountInBox int    `db:"amount_in_box" json:"amount_in_box"`
}

// amount of ordered good of certain type
type GoodOrdered struct {
	Good   Good `db:"good" json:"good"`
	Amount int  `db:"amount" json:"amount"`
}

// the order
type Order struct {
	Id                int           `db:"id" json:"id"`
	NumOrder          string        `db:"num_order" json:"num_order"`
	Contract          string        `db:"contract" json:"contract"`
	Run               int           `db:"run" json:"run"`
	Customer          string        `db:"customer" json:"customer"`
	OrderName         string        `db:"order_name" json:"order_name"`
	Address           string        `db:"address" json:"address"`
}

var connStr = "postgres://bbs_portal:JL84KdM_32@localhost/bbs_print_portal?sslmode=disable"


func GetAllOrdersForCompletion()([]OrdersModel,error){
	var result []OrdersModel

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("error establish connection: " + err.Error())
		return nil, err
	}
	defer db.Close()
	statementGetOrders := "select id, num_order, contract, run, customer, order_name, address from rosstat.rosstat_orders where completed = false;"
	rows, err := db.Query(statementGetOrders)
	if err != nil {
		log.Println("error query statementGetOrders - select all orders")
	}

	var index = 1
	for rows.Next(){

		order := Order{}
		rows.Scan(&order)

		boxes, err := getBoxesAmountForOrder(order.Id)
		if err != nil {
			log.Println("error get amount of box for order: " + err.Error())
		}

		pallets, err := getPalletsAmountForOrder(order.Id)
		if err != nil {
			log.Println("error get amount of pallets for order: " + err.Error())
		}

		smallBoxes, err := getSmallBoxesAmountForOrder(order.Id)
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
			AmountPallets: 0,
			AmountBoxes:   0,
			SubOrders:   []SubOrderModel{
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
	return result, nil
}
func getBoxesAmountForOrder(orderId int) (int,error){
	var boxes = 0
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("error establish connection: " + err.Error())
		return 0, err
	}
	defer db.Close()

	statement := "select count(id) from rosstat.boxes where pallet_id in " +
		"(select id from rosstat.pallets where order_id = " + strconv.Itoa(orderId) + ");"
	rows, err := db.Query(statement)
	if err != nil {
		log.Println("error get amount of boxes for order: " + err.Error())
	}
	for rows.Next(){
		rows.Scan(&boxes)
	}

	return boxes, nil
}

func getPalletsAmountForOrder(orderId int) (int,error){
	var pallets = 0
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("error establish connection: " + err.Error())
		return 0, err
	}
	defer db.Close()

	statement := "select count(id) from rosstat.pallets where order_id = " + strconv.Itoa(orderId) + ";"
	rows, err := db.Query(statement)
	if err != nil {
		log.Println("error get amount of pallets for order: " + err.Error())
	}
	for rows.Next(){
		rows.Scan(&pallets)
	}

	return pallets, nil
}

func getSmallBoxesAmountForOrder(orderId int) (int,error){
	var boxes = 0
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("error establish connection: " + err.Error())
		return 0, err
	}
	defer db.Close()

	statement := "select count(id) from rosstat.small_boxes where order_id = " + strconv.Itoa(orderId) + ";"
	rows, err := db.Query(statement)
	if err != nil {
		log.Println("error get amount of small boxes for order: " + err.Error())
	}
	for rows.Next(){
		rows.Scan(&boxes)
	}

	return boxes, nil
}





// Call it after button "combined boxes fully completed" pressed
// update data in rosstat_orders - subtract amount of goods that were completed
// we will not control every good itself, we believe that when operator clicked "all combined boxes completed - all pieces are packed"
func PutSmallOrderToDB(orderId int, boxIds []int, userName string) (int, error) {

	err := createSmallBoxesRecord(orderId, boxIds, userName)
	if err != nil {
		log.Println("Error create record in rosstat.small_boxes: " + err.Error())
		return 0, err
	}

	complectedGoods, err := GetAmountOfPiecesOfGood(orderId)
	if err != nil {
		log.Println("can't get amount of pieces in combined boxes: " + err.Error())
		return 0, err
	}

	err = updateDataInRosstatOrdersWithPieces(complectedGoods)

	if err != nil {
		log.Println("can't update order with small suborder: " + err.Error())
		return 0, err
	}

	return len(boxIds), nil
}

// TODO: update data in rosstat_orders - subtract amount of goods that were completed// Call it after button "Pallet is done" pressed and no error about small order completion, if this pallet is last
func CreatePallet(orderId, palletNum int, boxes []int, userName string) (int, int, int, error) {
	// id for pallet forms from order id and pallet number
	palletId, err := strconv.Atoi(strconv.Itoa(orderId) + strconv.Itoa(palletNum)) // result like: if orderId 1264 + palletNum 7 will be 12647
	if err != nil {
		log.Println("Cannot convert pallet id to int: " + err.Error())
		return 0, 0, 0, err
	}
	// TODO: put in db: insert into rosstat.pallets id, num, order_id values palletId, palletNum, orderId Error check!!!
	err = createBoxesRecord(palletId, boxes, userName)
	if err != nil {
		log.Println("Cannot convert pallet id to int: " + err.Error())
		return 0, 0, 0, err
	}

	return palletId, orderId, palletNum, nil
}

// Mostly this should be called when we start a new order. But can be called if order is not complete and closed by accidence or because of new shift
func GetLastPalletNumBuOrderId(orderId int) (int, error) {
	palletNum := 0
	// TODO: select max(num) from rosstat.pallets where oder_id = orderId and put it into palletNum  Error check!!!
	// if there is no pallets with this order_id - return 0, it means we have no pallets for this order yet
	return palletNum, nil
}

// list of all goods need to be collected in whole boxes. Called when "create pallet" button pressed in order completing window
func GetOrderListForWholeBoxes(orderId int) ([]GoodOrdered, error) {
	var result []GoodOrdered
	for i := 0; i < 26; i++ {
		tmp := GoodOrdered{}
		tmp.Good = GetProductByType(i + 1)
		amount := 0
		str := "select \"" + strconv.Itoa(i+1) + "\" from rosstat.rosstat_orders where id=" + strconv.Itoa(orderId) + ";"
		log.Println(str)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Println("error establish connection: " + err.Error())
			return nil, err
		}
		rows, err := db.Query(str)
		for rows.Next() {
			err = rows.Scan(&amount)
			if err != nil {
				log.Println("error scan row: " + err.Error())
				return nil, err
			}
		}
		tmp.Amount = GetWholeBoxesOfThisProduct(tmp.Good, amount) // here we get amount of boxes of this product!
		result = append(result, tmp)
	}
	return result, nil
}

// this list is for small part of order, for combined boxes(27 type) Called when small order clicked
func GetOrderListForPieces(orderId int) ([]GoodOrdered, error) {
	var result []GoodOrdered
	for i := 0; i < 26; i++ {
		tmp := GoodOrdered{}
		tmp.Good = GetProductByBoxID(i + 1)
		amount := 0
		// TODO: select strconv.Itoa(i+1) from rosstat.rosstat_orders where id = orderId. Put result into amount Error check!!!
		tmp.Amount = GetPiecesOfThisProduct(tmp.Good, amount) // here we get amount of pieces, not boxes of this product!
		result = append(result, tmp)
	}
	return result, nil
}

// ??
func GetOrderByNumber(orderNumber string) Order {
	order := Order{}
	// TODO: fill structure from db: select * from rosstat.rosstat_orders where num_order = orderNumber; for every good - call GetProductByID Error check!!!
	return order
}

// get amount of whole boxes of particular good
func GetWholeBoxesOfThisProduct(product Good, totalAmount int) int {
	return totalAmount / product.AmountInBox
}

// get amount of pieces of particular good for combined box
func GetPiecesOfThisProduct(product Good, totalAmount int) int {
	log.Println("total: " + strconv.Itoa(totalAmount) + ", Amount in a box: " + strconv.Itoa(product.AmountInBox))
	return totalAmount % product.AmountInBox
}

// returns amount of boxed for each type of good got from ids array
func GetAmountOfBoxesOfGoodToSubtractFromDB(boxIds []int) ([27]int, error) {
	var result [27]int
	for i := 0; i < len(boxIds); i++ {
		goodType, err := strconv.Atoi(GetProductByBoxID(boxIds[i]).Num)
		if err != nil {
			log.Println("Error try to convert good type to int: " + err.Error())
			return [27]int{}, err
		}
		boxIds[goodType]++
	}
	return result, nil
}

// returns amount of pieces for each type of good got from ids array
func GetAmountOfPiecesOfGood(orderId int) ([26]int, error) {
	var result [26]int
	statement := "select "
	for i := 1; i < 27; i++ {
		statement += strconv.Itoa(i) + ", "
	}
	statement = strings.TrimRight(statement, ",")
	statement += " from rosstat.rosstat_orders where id = "
	statement += strconv.Itoa(orderId)
	// TODO: exec statement(notice if in cell null - it means 0), put to result
	for i := 0; i < 26; i++ {
		// we got total amount of every good here, but we need only pieces, so update result
		result[i] = GetPiecesOfThisProduct(GetProductByBoxID(i+1), result[i])
	}
	return result, nil
}

func GetProductByBoxID(id int) Good {
	result := Good{
	}

	if id > 200100001 && id < 200200000 {
		result.Name = "Форма № 1. Записная книжка переписчика (является приложением к Инструкции)"
		result.Run = 476596
		result.AmountInBox = 20
		result.Num = "1"

	} else if id > 200200001 && id < 200300000 {
		result.Name = "Форма № 2. Записная книжка контролера полевого уровня"
		result.Run = 65357
		result.AmountInBox = 50
		result.Num = "2"

	} else if id > 200300001 && id < 200400000 {
		result.Name = "Форма № 3. Записная книжка уполномоченного по вопросам переписи"
		result.Run = 6023
		result.AmountInBox = 50
		result.Num = "3"

	} else if id > 200400001 && id < 200500000 {
		result.Name = "Форма № 4. Сводная ведомость по переписному участку"
		result.Run = 57930
		result.AmountInBox = 1000
		result.Num = "4"

	} else if id > 200500001 && id < 200600000 {
		result.Name = "Форма № 5. Сводная ведомость по городскому округу, муниципальному району/ округу"
		result.Run = 10459
		result.AmountInBox = 1000
		result.Num = "5"

	} else if id > 200600001 && id < 200700000 {
		result.Name = "Форма № 6. Сводка итогов переписи населения по городскому округу, муниципальному району/округу"
		result.Run = 61459
		result.AmountInBox = 500
		result.Num = "6"

	} else if id > 200700001 && id < 200800000 {
		result.Name = "Форма № 7. Информационные листовки (к лицам, которых трудно застать дома)"
		result.Run = 28419540
		result.AmountInBox = 2000
		result.Num = "7"

	} else if id > 200800001 && id < 200900000 {
		result.Name = "Форма № 9. Ярлык в портфель переписчика"
		result.Run = 18812
		result.AmountInBox = 8000
		result.Num = "8"

	} else if id > 200900001 && id < 201000000 {
		result.Name = "Форма № 10. Карточка для респондентов"
		result.Run = 392150
		result.AmountInBox = 2000
		result.Num = "9"

	} else if id > 201000001 && id < 201100000 {
		result.Name = "Форма Обложка. Обложка на переписные документы"
		result.Run = 2287448
		result.AmountInBox = 500
		result.Num = "10"

	} else if id > 201100001 && id < 201200000 {
		result.Name = "Форма С. Список лиц"
		result.Run = 2287448
		result.AmountInBox = 1000
		result.Num = "11"

	} else if id > 201200001 && id < 201300000 {
		result.Name = "Форма КС. Список лиц для контроля за заполнением переписных листов"
		result.Run = 790907
		result.AmountInBox = 2000
		result.Num = "12"

	} else if id > 201300001 && id < 201400000 {
		result.Name = "Форма СПР. Справка о прохождении переписи"
		result.Run = 10136033
		result.AmountInBox = 8000
		result.Num = "13"

	} else if id > 201400001 && id < 201500000 {
		result.Name = "Инструкция о порядке подготовки материалов Всероссийской переписи населения 2020 года к обработке"
		result.Run = 1652
		result.AmountInBox = 40
		result.Num = "14"

	} else if id > 201500001 && id < 201600000 {
		result.Name = "Тесты для обучения переписного персонала"
		result.Run = 495982
		result.AmountInBox = 100
		result.Num = "15"

	} else if id > 201600001 && id < 201700000 {
		result.Name = "Указатели для переписных участков"
		result.Run = 32689
		result.AmountInBox = 500
		result.Num = "16"

	} else if id > 201700001 && id < 201800000 {
		result.Name = "Форма Л. Переписной лист (обучение, чистые формы)"
		result.Run = 2082395
		result.AmountInBox = 1000
		result.Num = "17"

	} else if id > 201800001 && id < 201900000 {
		result.Name = "Форма П. Переписной лист (обучение, чистые формы)"
		result.Run = 832958
		result.AmountInBox = 1000
		result.Num = "18"

	} else if id > 201900001 && id < 202000000 {
		result.Name = "Форма В. Переписной лист (обучение, чистые формы)"
		result.Run = 416479
		result.AmountInBox = 1000
		result.Num = "19"

	} else if id > 202000001 && id < 202100000 {
		result.Name = "Форма Н. Сопроводительный бланк (обучение, чистые формы)"
		result.Run = 416479
		result.AmountInBox = 1000
		result.Num = "20"

	} else if id > 202100001 && id < 202200000 {
		result.Name = "Форма Обложка. Обложка на переписные документы (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 500
		result.Num = "21"

	} else if id > 202200001 && id < 202300000 {
		result.Name = "Форма С. Список лиц (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Num = "22"

	} else if id > 202300001 && id < 202400000 {
		result.Name = "Форма Л. Переписной лист (обучение, заполненные формы)"
		result.Run = 8260
		result.AmountInBox = 1000
		result.Num = "23"

	} else if id > 202400001 && id < 202500000 {
		result.Name = "Форма П. Переписной лист (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Num = "24"

	} else if id > 202500001 && id < 202600000 {
		result.Name = "Форма В. Переписной лист (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Num = "25"

	} else if id > 202600001 && id < 202700000 {
		result.Name = "Форма Н. Сопроводительный бланк (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Num = "26"

	} else if id > 202700001 && id < 202800000 {
		result.Name = "Сборный короб"
		result.Run = 0
		result.AmountInBox = 0
		result.Num = "27"

	}

	return result
}

func GetProductByType(t int) Good {
	result := Good{
	}

	switch t {
	case 1:
		result.Name = "Форма № 1. Записная книжка переписчика (является приложением к Инструкции)"
		result.Run = 476596
		result.AmountInBox = 20
		result.Num = "1"
		break
	case 2:
		result.Name = "Форма № 2. Записная книжка контролера полевого уровня"
		result.Run = 65357
		result.AmountInBox = 50
		result.Num = "2"
		break
	case 3:
		result.Name = "Форма № 3. Записная книжка уполномоченного по вопросам переписи"
		result.Run = 6023
		result.AmountInBox = 50
		result.Num = "3"
		break
	case 4:
		result.Name = "Форма № 4. Сводная ведомость по переписному участку"
		result.Run = 57930
		result.AmountInBox = 1000
		result.Num = "4"
		break
	case 5:
		result.Name = "Форма № 5. Сводная ведомость по городскому округу, муниципальному району/ округу"
		result.Run = 10459
		result.AmountInBox = 1000
		result.Num = "5"
		break
	case 6:
		result.Name = "Форма № 6. Сводка итогов переписи населения по городскому округу, муниципальному району/округу"
		result.Run = 61459
		result.AmountInBox = 500
		result.Num = "6"
		break
	case 7:
		result.Name = "Форма № 7. Информационные листовки (к лицам, которых трудно застать дома)"
		result.Run = 28419540
		result.AmountInBox = 2000
		result.Num = "7"
		break
	case 8:
		result.Name = "Форма № 9. Ярлык в портфель переписчика"
		result.Run = 18812
		result.AmountInBox = 8000
		result.Num = "8"
		break
	case 9:
		result.Name = "Форма № 10. Карточка для респондентов"
		result.Run = 392150
		result.AmountInBox = 2000
		result.Num = "9"
		break
	case 10:
		result.Name = "Форма Обложка. Обложка на переписные документы"
		result.Run = 2287448
		result.AmountInBox = 500
		result.Num = "10"
		break
	case 11:
		result.Name = "Форма С. Список лиц"
		result.Run = 2287448
		result.AmountInBox = 1000
		result.Num = "11"
		break
	case 12:
		result.Name = "Форма КС. Список лиц для контроля за заполнением переписных листов"
		result.Run = 790907
		result.AmountInBox = 2000
		result.Num = "12"

		break
	case 13:
		result.Name = "Форма СПР. Справка о прохождении переписи"
		result.Run = 10136033
		result.AmountInBox = 8000
		result.Num = "13"
		break
	case 14:
		result.Name = "Инструкция о порядке подготовки материалов Всероссийской переписи населения 2020 года к обработке"
		result.Run = 1652
		result.AmountInBox = 40
		result.Num = "14"
		break
	case 15:
		result.Name = "Тесты для обучения переписного персонала"
		result.Run = 495982
		result.AmountInBox = 100
		result.Num = "15"
		break
	case 16:
		result.Name = "Указатели для переписных участков"
		result.Run = 32689
		result.AmountInBox = 500
		result.Num = "16"
		break
	case 17:
		result.Name = "Форма Л. Переписной лист (обучение, чистые формы)"
		result.Run = 2082395
		result.AmountInBox = 1000
		result.Num = "17"
		break
	case 18:
		result.Name = "Форма П. Переписной лист (обучение, чистые формы)"
		result.Run = 832958
		result.AmountInBox = 1000
		result.Num = "18"
		break
	case 19:
		result.Name = "Форма В. Переписной лист (обучение, чистые формы)"
		result.Run = 416479
		result.AmountInBox = 1000
		result.Num = "19"
		break
	case 20:
		result.Name = "Форма Н. Сопроводительный бланк (обучение, чистые формы)"
		result.Run = 416479
		result.AmountInBox = 1000
		result.Num = "20"
		break
	case 21:
		result.Name = "Форма Обложка. Обложка на переписные документы (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 500
		result.Num = "21"
		break
	case 22:
		result.Name = "Форма С. Список лиц (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Num = "22"
		break
	case 23:
		result.Name = "Форма Л. Переписной лист (обучение, заполненные формы)"
		result.Run = 8260
		result.AmountInBox = 1000
		result.Num = "23"
		break
	case 24:
		result.Name = "Форма П. Переписной лист (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Num = "24"
		break
	case 25:
		result.Name = "Форма В. Переписной лист (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Num = "25"
		break
	case 26:
		result.Name = "Форма Н. Сопроводительный бланк (обучение, заполненные формы)"
		result.Run = 1652
		result.AmountInBox = 1000
		result.Num = "26"
		break
	case 27:
		result.Name = "Сборный короб"
		result.Run = 0
		result.AmountInBox = 0
		result.Num = "27"
		break

	}

	return result
}

// put data in rosstat.boxes
func createBoxesRecord(palletId int, boxes []int, userName string) error {
	statement := "insert into rosstat.boxes values " // rosstat.boxes columns: id, pallet_id, user_name
	palletIdStr := strconv.Itoa(palletId)
	for i := 0; i < len(boxes); i++ {
		statement += "(" + strconv.Itoa(boxes[i]) + ", " + palletIdStr + ", " + userName + "),"
	}
	statement = strings.TrimRight(statement, ",")
	statement += ";"
	// TODO: put in db statement Error check!!!
	return nil
}

// put data in rosstat.small_boxes
func createSmallBoxesRecord(orderId int, boxIds []int, userName string) error {
	statementInsert := "insert into rosstat.small_boxes values "
	for i := 0; i < len(boxIds); i++ {
		statementInsert += "(" + strconv.Itoa(boxIds[i]) + ", " + strconv.Itoa(orderId) + ", " + userName + "),"
	}
	statementInsert = strings.TrimRight(statementInsert, ",")
	statementInsert += ";"
	// TODO: put in db statementInsert Error check!!!
	return nil
}

// subtract collected boxes from order
func updateDataInRosstatOrdersWithBoxes() error {
	return nil
}

// subtract collected pieces of goods from order
func updateDataInRosstatOrdersWithPieces(complectedGoods [26]int) error {
	statementUpdate := "update rosstat.rosstat_orders set "
	for i := 1; i < 27; i++ {
		statementUpdate += strconv.Itoa(i) + "=" + strconv.Itoa(i) + "- " + strconv.Itoa(complectedGoods[i-1]) + ", "
	}
	statementUpdate = strings.TrimRight(statementUpdate, ",")
	statementUpdate += ";"
	// TODO: run statement Error check!!!
	return nil
}
