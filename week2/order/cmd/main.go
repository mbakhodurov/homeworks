package main

import (
	"context"
	"net/http"
	"time"

	order_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/payment/v1"
)

const (
	inventoryAddress = "localhost:50052"
	paymentAddress   = "localhost:50051"

	httpPort = "8086"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type Handler struct {
	inventory inventory_v1.InventoryServiceClient
	payment   payment_v1.PaymentServiceClient
}

func NewHandler(inventory inventory_v1.InventoryServiceClient, payment payment_v1.PaymentServiceClient) *Handler {
	return &Handler{
		inventory: inventory,
		payment:   payment,
	}
}

func (h *Handler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: order_v1.GenericError{
			Code:    order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
	}
}

func CreateOrder()

func main() {

}
