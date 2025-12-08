package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week1/order/pkg/models"
	order_v1 "github.com/mbakhodurov/homeworks/week1/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/mbakhodurov/homeworks/week1/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	inventoryAddress = "localhost:50052"

	httpPort = "8086"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderHandler struct {
	storage   *models.OrderStorage
	inventory inventory_v1.InventoryServiceClient
}

func NewOrderHandler(storage *models.OrderStorage, inventory inventory_v1.InventoryServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:   storage,
		inventory: inventory,
	}
}

func (h *OrderHandler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: order_v1.GenericError{
			Code:    order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
	}
}

func (h *OrderHandler) GetOrderById(ctx context.Context, params order_v1.GetOrderByIdParams) (order_v1.GetOrderByIdRes, error) {
	orderUUID := params.OrderUUID

	order, err := h.storage.GetOrderByUUID(orderUUID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Заказ не найден",
			}, nil
		}
		return nil, err
	}
	dto := order_v1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: float32(order.TotalPrice),
		Status:     order_v1.OrderStatus(order.Status),
	}
	if order.TransactionUUID != nil {
		dto.TransactionUUID = order_v1.NewOptNilString(*order.TransactionUUID)
	}
	if order.PaymentMethod != nil {
		dto.PaymentMethod = &order_v1.NilOrderDtoPaymentMethod{Value: order_v1.OrderDtoPaymentMethod(*order.PaymentMethod)}
	}
	return &order_v1.GetOrderResponse{
		OrderDto: dto,
	}, nil
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	if len(req.PartUuids) == 0 {
		return &order_v1.BadRequestError{
			Code:    400,
			Message: "required PartUuids",
		}, nil
	}

	parts, err := h.inventory.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: req.PartUuids,
		},
	})

	// if err != nil {
	// 	return &order_v1.BadRequestError{
	// 		Code:    400,
	// 		Message: "some parts not found",
	// 	}, nil
	// }

	if err != nil || int(parts.GetTotalCount()) != len(req.PartUuids) {
		if err != nil {
			return &order_v1.BadRequestError{
				Code:    400,
				Message: "some parts not found",
			}, nil
		}
		return &order_v1.BadRequestError{
			Code:    400,
			Message: "some parts not found",
		}, nil
	}
	var total_price float64 = 0
	for _, v := range parts.Part {
		total_price += v.Info.Price
	}

	orderNew := models.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   req.GetUserUUID(),
		PartUUIDs:  req.GetPartUuids(),
		TotalPrice: total_price,
		Status:     models.StatusPending,
	}
	h.storage.CreateOrder(&orderNew)

	return &order_v1.CreateOrderResponse{
		OrderUUID:  orderNew.OrderUUID,
		TotalPrice: float32(orderNew.TotalPrice),
	}, nil
}

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryConn)

	orderHandler := NewOrderHandler(models.NewStorage(), inventoryClient)

	orderServer, err := order_v1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
