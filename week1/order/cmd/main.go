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
	payment_v1 "github.com/mbakhodurov/homeworks/week1/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	storage   *models.OrderStorage
	inventory inventory_v1.InventoryServiceClient
	payment   payment_v1.PaymentServiceClient
}

func NewHandler(storage *models.OrderStorage, inventory inventory_v1.InventoryServiceClient, payment payment_v1.PaymentServiceClient) *Handler {
	return &Handler{
		storage:   storage,
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

func (h *Handler) GetOrderById(ctx context.Context, params order_v1.GetOrderByIdParams) (order_v1.GetOrderByIdRes, error) {
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

func (h *Handler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
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

func (h *Handler) PaymentOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PaymentOrderParams) (order_v1.PaymentOrderRes, error) {
	order, err := h.storage.GetOrderByUUID(params.OrderUUID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Заказ не найден",
			}, nil
		}
		return nil, err
	}

	if order.Status != models.StatusPending {
		return &order_v1.ConflictError{
			Code:    409,
			Message: "Order not in pending status",
		}, nil
	}

	var paymentMethod payment_v1.PaymentMethod
	switch req.PaymentMethod {
	case "CARD":
		paymentMethod = payment_v1.PaymentMethod_CARD
	case "SBP":
		paymentMethod = payment_v1.PaymentMethod_SBP
	case "CREDIT_CARD":
		paymentMethod = payment_v1.PaymentMethod_CREDIT_CARD
	case "INVESTOR_MONEY":
		paymentMethod = payment_v1.PaymentMethod_INVESTOR_MONEY
	default:
		return &order_v1.InternalServerError{Code: 500, Message: "Invalid payment method"}, nil
	}

	payresp, err := h.payment.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return &order_v1.InternalServerError{Code: 500, Message: "Err:=" + err.Error()}, nil
	}
	order.Status = models.StatusPaid
	order.PaymentMethod = (*models.PaymentMethod)(&req.PaymentMethod)
	return &order_v1.PayOrderResponse{
		OrderUUID: order_v1.OptString{Value: payresp.TransactionUuid},
	}, nil
}

func (h *Handler) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	order, err := h.storage.GetOrderByUUID(params.OrderUUID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Заказ не найден",
			}, nil
		}
		return nil, err
	}

	if order.Status == models.StatusPaid {
		return &order_v1.ConflictError{
			Code:    409,
			Message: "Order has already paid",
		}, nil
	}

	order.Status = models.StatusCancelled
	return &order_v1.CancelOrderResponse{}, nil
}

func (h *Handler) GetAllOrders(ctx context.Context) (order_v1.GetAllOrdersRes, error) {
	orders, err := h.storage.GetAllOrder()
	if err != nil {
		return &order_v1.InternalServerError{
			Code:    500,
			Message: err.Error(),
		}, nil
	}
	if len(orders) == 0 {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Orders not found",
		}, nil
	}

	dtoList := make([]order_v1.OrderDto, 0, len(orders))
	for _, o := range orders {
		dto := order_v1.OrderDto{
			OrderUUID:  o.OrderUUID,
			UserUUID:   o.UserUUID,
			PartUuids:  o.PartUUIDs,
			TotalPrice: float32(o.TotalPrice),
			Status:     order_v1.OrderStatus(o.Status),
		}
		if o.TransactionUUID != nil {
			dto.TransactionUUID = order_v1.NewOptNilString(*o.TransactionUUID)
		}

		// payment_method — optional
		if o.PaymentMethod != nil {
			dto.PaymentMethod = &order_v1.NilOrderDtoPaymentMethod{Value: order_v1.OrderDtoPaymentMethod(*o.PaymentMethod)}
		}
		dtoList = append(dtoList, dto)
	}
	return &order_v1.GetAllOrderResponse{
		OrderDto: dtoList,
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

	paymentConn, err := grpc.NewClient(
		paymentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()
	paymentClient := payment_v1.NewPaymentServiceClient(paymentConn)

	orderHandler := NewHandler(models.NewStorage(), inventoryClient, paymentClient)

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
