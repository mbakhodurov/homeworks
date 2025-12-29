package main

// import (
// 	"context"
// 	"errors"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/google/uuid"
// 	"github.com/mbakhodurov/homeworks/week2/order/pkg/models"
// 	order_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/openapi/order/v1"
// 	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
// 	payment_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/payment/v1"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// const (
// 	inventoryAddress = "localhost:50052"
// 	paymentAddress   = "localhost:50051"

// 	httpPort = "8086"
// 	// Таймауты для HTTP-сервера
// 	readHeaderTimeout = 5 * time.Second
// 	shutdownTimeout   = 10 * time.Second
// )

// type Handler struct {
// 	storage   *models.OrderStorage
// 	inventory inventory_v1.InventoryServiceClient
// 	payment   payment_v1.PaymentServiceClient
// }

// func NewHandler(storage *models.OrderStorage, inventory inventory_v1.InventoryServiceClient, payment payment_v1.PaymentServiceClient) *Handler {
// 	return &Handler{
// 		inventory: inventory,
// 		payment:   payment,
// 		storage:   storage,
// 	}
// }

// func (h *Handler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
// 	if len(req.PartUuids) == 0 {
// 		return &order_v1.BadRequestError{
// 			Code:    400,
// 			Message: "required PartUuids",
// 		}, nil
// 	}
// 	parts, err := h.inventory.ListParts(ctx, &inventory_v1.ListPartsRequest{
// 		Filter: &inventory_v1.PartsFilter{
// 			Uuids: req.GetPartUuids(),
// 		},
// 	})

// 	if err != nil || int(parts.GetTotalCount()) != len(req.PartUuids) {
// 		if err != nil {
// 			return &order_v1.BadRequestError{
// 				Code:    400,
// 				Message: "some parts not found",
// 			}, nil
// 		}
// 		return &order_v1.BadRequestError{
// 			Code:    400,
// 			Message: "some parts not found",
// 		}, nil
// 	}

// 	total_price := 0
// 	for _, v := range parts.Part {
// 		total_price += int(v.Info.Price)
// 	}
// 	order_new := models.Order{
// 		OrderUUID:  uuid.NewString(),
// 		UserUUID:   req.GetUserUUID(),
// 		PartUUIDS:  req.GetPartUuids(),
// 		TotalPrice: float64(total_price),
// 		Status:     models.StatusPending,
// 	}
// 	h.storage.Create(&order_new)
// 	return &order_v1.CreateOrderResponse{
// 		OrderUUID:  order_new.OrderUUID,
// 		TotalPrice: float32(order_new.TotalPrice),
// 	}, nil
// }

// func (h *Handler) GetAllOrders(ctx context.Context) (order_v1.GetAllOrdersRes, error) {
// 	orders, err := h.storage.GetAll()
// 	if err != nil {
// 		if errors.Is(err, models.ErrNotFound) {
// 			return &order_v1.NotFoundError{
// 				Code:    404,
// 				Message: "Заказов пока что нету",
// 			}, nil
// 		}
// 		return &order_v1.InternalServerError{
// 			Code:    500,
// 			Message: err.Error(),
// 		}, nil
// 	}
// 	orderDtoList := make([]order_v1.OrderDto, 0, len(orders))

// 	for _, v := range orders {
// 		dto := order_v1.OrderDto{
// 			OrderUUID:  v.OrderUUID,
// 			UserUUID:   v.UserUUID,
// 			PartUuids:  v.PartUUIDS,
// 			TotalPrice: float32(v.TotalPrice),
// 			Status:     order_v1.OrderStatus(v.Status),
// 		}

// 		if v.TransactionUUID != nil {
// 			dto.TransactionUUID = order_v1.NewOptNilString(*v.TransactionUUID)
// 		}

// 		if v.PaymentMethod != nil {
// 			dto.PaymentMethod = &order_v1.NilOrderDtoPaymentMethod{
// 				Value: order_v1.OrderDtoPaymentMethod(*v.PaymentMethod),
// 			}
// 		}
// 		orderDtoList = append(orderDtoList, dto)
// 	}

// 	return &order_v1.GetAllOrderResponse{
// 		OrderDto:   orderDtoList,
// 		TotalCount: float64(len(orderDtoList)),
// 	}, nil
// }

// func (h *Handler) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
// 	order, err := h.storage.Get(params.OrderUUID)
// 	if err != nil {
// 		if errors.Is(err, models.ErrNotFound) {
// 			return &order_v1.NotFoundError{
// 				Code:    404,
// 				Message: "Заказ не найден",
// 			}, nil
// 		}
// 		return nil, err
// 	}

// 	if order.Status == models.StatusPaid {
// 		return &order_v1.ConflictError{
// 			Code:    409,
// 			Message: "Order has already paid",
// 		}, nil
// 	}

// 	order.Status = models.StatusCancelled
// 	return &order_v1.CancelOrderResponse{}, nil
// }

// func (h *Handler) PaymentOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PaymentOrderParams) (order_v1.PaymentOrderRes, error) {
// 	order, err := h.storage.Get(params.OrderUUID)
// 	if err != nil {
// 		if errors.Is(err, models.ErrNotFound) {
// 			return &order_v1.NotFoundError{
// 				Code:    404,
// 				Message: "Заказ не найден",
// 			}, nil
// 		}
// 		return nil, err
// 	}
// 	if order.Status != models.StatusPending {
// 		return &order_v1.ConflictError{
// 			Code:    409,
// 			Message: "Order not in pending status",
// 		}, nil
// 	}

// 	var paymentMethod payment_v1.PaymentMethod
// 	switch req.PaymentMethod {
// 	case "CARD":
// 		paymentMethod = payment_v1.PaymentMethod_CARD
// 	case "SBP":
// 		paymentMethod = payment_v1.PaymentMethod_SBP
// 	case "CREDIT_CARD":
// 		paymentMethod = payment_v1.PaymentMethod_CREDIT_CARD
// 	case "INVESTOR_MONEY":
// 		paymentMethod = payment_v1.PaymentMethod_INVESTOR_MONEY
// 	default:
// 		return &order_v1.InternalServerError{
// 			Code:    500,
// 			Message: "Invalid payment method",
// 		}, nil
// 	}

// 	payresp, err := h.payment.PayOrder(ctx, &payment_v1.PayOrderRequest{
// 		OrderUuid:     order.OrderUUID,
// 		UserUuid:      order.UserUUID,
// 		PaymentMethod: paymentMethod,
// 	})
// 	if err != nil {
// 		return &order_v1.InternalServerError{Code: 500, Message: "Err:=" + err.Error()}, nil
// 	}

// 	order.Status = models.StatusPaid
// 	order.PaymentMethod = (*models.PaymentMethod)(&req.PaymentMethod)

// 	return &order_v1.PayOrderResponse{
// 		OrderUUID: order_v1.OptString{Value: payresp.TransactionUuid},
// 	}, nil
// }

// func (h *Handler) GetOrderByUUID(ctx context.Context, params order_v1.GetOrderByUUIDParams) (order_v1.GetOrderByUUIDRes, error) {
// 	order, err := h.storage.Get(params.OrderUUID)
// 	if err != nil {
// 		if errors.Is(err, models.ErrNotFound) {
// 			return &order_v1.NotFoundError{
// 				Code:    404,
// 				Message: "Заказ не найден",
// 			}, nil
// 		}
// 		return nil, err
// 	}
// 	order_dto := order_v1.OrderDto{
// 		OrderUUID:  order.OrderUUID,
// 		UserUUID:   order.UserUUID,
// 		PartUuids:  order.PartUUIDS,
// 		TotalPrice: float32(order.TotalPrice),
// 	}

// 	if order.TransactionUUID != nil {
// 		order_dto.TransactionUUID = order_v1.NewOptNilString(*order.TransactionUUID)
// 	}
// 	if order.PaymentMethod != nil {
// 		order_dto.PaymentMethod = &order_v1.NilOrderDtoPaymentMethod{Value: order_v1.OrderDtoPaymentMethod(*order.PaymentMethod)}
// 	}

// 	return &order_v1.GetOrderResponse{
// 		OrderDto: order_dto,
// 	}, nil
// }

// func (h *Handler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
// 	return &order_v1.GenericErrorStatusCode{
// 		StatusCode: 500,
// 		Response: order_v1.GenericError{
// 			Code:    order_v1.NewOptInt(http.StatusInternalServerError),
// 			Message: order_v1.NewOptString(err.Error()),
// 		},
// 	}
// }

// func main() {
// 	inventoryConn, err := grpc.NewClient(
// 		inventoryAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		log.Printf("failed to connect: %v\n", err)
// 		return
// 	}
// 	defer func() {
// 		if cerr := inventoryConn.Close(); cerr != nil {
// 			log.Printf("failed to close connect: %v", cerr)
// 		}
// 	}()
// 	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryConn)

// 	paymentConn, err := grpc.NewClient(
// 		paymentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		log.Printf("failed to connect: %v\n", err)
// 		return
// 	}
// 	defer func() {
// 		if cerr := paymentConn.Close(); cerr != nil {
// 			log.Printf("failed to close connect: %v", cerr)
// 		}
// 	}()
// 	paymentClient := payment_v1.NewPaymentServiceClient(paymentConn)

// 	orderHandler := NewHandler(models.NewOrderStorage(), inventoryClient, paymentClient)

// 	orderServer, err := order_v1.NewServer(orderHandler)
// 	if err != nil {
// 		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
// 	}
// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)
// 	r.Use(middleware.Timeout(10 * time.Second))

// 	r.Mount("/", orderServer)

// 	server := &http.Server{
// 		Addr:              net.JoinHostPort("localhost", httpPort),
// 		Handler:           r,
// 		ReadHeaderTimeout: readHeaderTimeout,
// 	}

// 	go func() {
// 		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
// 		err = server.ListenAndServe()
// 		if err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
// 		}
// 	}()

// 	// Graceful shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	log.Println("🛑 Завершение работы сервера...")

// 	// Создаем контекст с таймаутом для остановки сервера
// 	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
// 	defer cancel()

// 	err = server.Shutdown(ctx)
// 	if err != nil {
// 		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
// 	}

// 	log.Println("✅ Сервер остановлен")
// }
