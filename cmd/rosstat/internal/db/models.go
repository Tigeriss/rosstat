package db

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
