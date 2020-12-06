package db

// internal models

// good itself. hardcoded because of problems with id
type Good struct {
	Type        int    `db:"num" json:"num"`
	Name        string `db:"name" json:"name"`
	Run         int    `db:"run" json:"run"`
	AmountInBox int    `db:"amount_in_box" json:"amount_in_box"`
	FirstID     int    `db:"first_id" json:"first_id"`
	LastID      int    `db:"last_id" json:"last_id"`
}

type PalletRegistryGoodsData struct {
	good Good
	boxes int
}

// the order
type Order struct {
	Id        int    `db:"id" json:"id"`
	NumOrder  string `db:"num_order" json:"num_order"`
	Contract  string `db:"contract" json:"contract"`
	Run       int    `db:"run" json:"run"`
	Customer  string `db:"customer" json:"customer"`
	OrderName string `db:"order_name" json:"order_name"`
	Address   string `db:"address" json:"address"`
}

type PalletInfo struct {
	barcode string
	palletNum int
	boxes int
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
	SubOrders     []SubOrderModel `json:"sub_orders,omitempty"`
}

// GET /orders/big/build - used by /orders/big page

type BigOrdersModel struct {
	Type     int    `json:"type"`
	FormName string `json:"form_name"`
	Total    int    `json:"total"`
	Built    int    `json:"built"`
}

type FinishSmallOrderModel struct {
	Boxes []string `json:"boxes"`
}

// GET /orders/big/pallet/:id - from /order/pallet page

type BigPalletModel struct {
	PalletNum int              `json:"pallet_num"`
	Types     []BigOrdersModel `json:"types"`
}

type BigPalletBarcodeModel struct {
	Success bool   `json:"success"`
	Type    int    `json:"type"`
	Error   string `json:"error"`
}

type BigPalletFinishRequestModel struct {
	PalletNum int      `json:"pallet_num"`
	Barcodes  []string `json:"barcodes"`
}

type BigPalletFinishResponseModel struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	LastPallet bool   `json:"last_pallet"`
}

// GET /shipment/ready

type ShipmentPalletModel struct {
	Num         int    `json:"num"`
	PalletNum   int    `json:"pallet_num"`
	Barcode     string `json:"barcode"`
	AmountBoxes int    `json:"amount_boxes"`
}

type PrintPalletRegisterModel struct {
	NumPP    int    `json:"num_pp"`
	Position string `json:"position"`
	Amount   int    `json:"amount"`
	Boxes    int    `json:"boxes"`
}

type PrintPalletModel struct {
	OrderCaption   string                     `json:"order_caption"`
	Address        string                     `json:"address"`
	Provider       string                     `json:"provider"`
	ContractNumber string                     `json:"contract_number"`
	Barcode        string                     `json:"barcode"`
	Register       []PrintPalletRegisterModel `json:"register"`
}
