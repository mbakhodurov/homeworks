package model

import "time"

type OrderUpdateInfo struct {
	Payment_method *PaymentMethod `json:"payment_method"`
	Updated_at     *time.Time     `json:"updated_at"`
	Deleted_at     *time.Time     `json:"deleted_at"`
	Status         *OrderStatus   `json:"status"`
}

type Order struct {
	Order_uuid       string         `json:"order_uuid"`
	User_uuid        string         `json:"user_uuid"`
	Part_uuids       []string       `json:"part_uuids"`
	Total_price      float64        `json:"total_price"`
	Transaction_uuid *string        `json:"transaction_uuid "`
	Payment_method   *PaymentMethod `json:"payment_method"`
	Status           OrderStatus    `json:"status"`

	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Deleted_at time.Time `json:"deleted_at"`
}

type OrderStatus int

const (
	StatusPendingPayment OrderStatus = iota
	StatusPaid
	StatusCancelled
)

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)
