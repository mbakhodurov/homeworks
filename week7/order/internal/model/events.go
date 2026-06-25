package model

type OrderPaid struct {
	EventUUID       string `json:"event_uuid"`
	OrderUUID       string `json:"order_uuid"`
	UserUUID        string `json:"user_uuid"`
	PaymentMethod   string `json:"payment_method"`
	TransactionUUID string `json:"transaction_uuid"`
}

type ShipAssembled struct {
	EventUUID    string // Уникальный идентификатор события (для идемпотентности)
	OrderUUID    string // Идентификатор собранного заказа
	UserUUID     string // Идентификатор пользователя
	BuildTimeSec int64  // Время (в секундах), потраченное на сборку корабля
}
