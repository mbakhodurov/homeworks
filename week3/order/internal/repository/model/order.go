package model

import "time"

type OrderUpdateInfo struct {
	Status           *OrderStatus   `json:"status"`
	Transaction_uuid *string        `json:"transaction_uuid"`
	Payment_method   *PaymentMethod `json:"payment_method"`
	Updated_at       *time.Time     `json:"updated_at"`
	Deleted_at       *time.Time     `json:"deleted_at"`
}

type Order struct {
	OrderUUID       string         `json:"order_uuid"`
	UserUUID        string         `json:"user_uuid"`
	PartUUIDs       []string       `json:"part_uuids"`
	TotalPrice      float64        `json:"total_price"`
	TransactionUUID *string        `json:"transaction_uuid"`
	PaymentMethod   *PaymentMethod `json:"payment_method"`
	Status          OrderStatus    `json:"status"`
	CreatedAt       *time.Time     `json:"created_at"`
	Updated_at      *time.Time     `json:"updated_at"`
	Deleted_at      *time.Time     `json:"deleted_at"`
}

type OrderStatus int

const (
	StatusPendingPayment OrderStatus = iota
	StatusPaid
	StatusCancelled
	StatusCompleted
)

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)
