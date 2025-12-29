package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/mbakhodurov/homeworks/week2/inventory/internal/api/inventory/v1"
	"github.com/mbakhodurov/homeworks/week2/inventory/internal/interceptor"
	"github.com/mbakhodurov/homeworks/week2/inventory/internal/repository/inventory"
	service "github.com/mbakhodurov/homeworks/week2/inventory/internal/service/inventory"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort = 50052
	httpPort = 8082
)

type InventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer
	mu sync.RWMutex

	inventory map[string]*inventory_v1.Part
}

func (i *InventoryService) Update(ctx context.Context, req *inventory_v1.UpdateRequest) (*emptypb.Empty, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	if req.UpdateInfo == nil {
		return nil, status.Error(codes.InvalidArgument, "update_info is required")
	}

	part, ok := i.inventory[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", req.GetUuid())
	}

	ui := req.GetUpdateInfo()
	// info := part.Info

	if ui.Name != nil {
		part.Info.Name = ui.GetName().Value
	}

	if ui.Description != nil {
		part.Info.Description = ui.GetDescription().Value
	}

	if ui.Price != nil {
		part.Info.Price = ui.GetPrice().Value
	}

	if ui.StockQuantity != nil {
		part.Info.StockQuantity = ui.GetStockQuantity().Value
	}

	if ui.Dimensions != nil {
		part.Info.Dimensions = ui.GetDimensions()
	}

	if ui.Manufacturer != nil {
		part.Info.Manufacturer = ui.GetManufacturer()
	}

	// updated_at
	part.UpdatedAt = timestamppb.New(time.Now())

	return &emptypb.Empty{}, nil
}

func (i *InventoryService) ListParts(_ context.Context, rq *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	filter := rq.GetFilter()
	result := []*inventory_v1.Part{}

	for _, part := range i.inventory {
		if matchPart(part, filter) {
			result = append(result, part)
		}
	}

	return &inventory_v1.ListPartsResponse{
		Part:       result,
		TotalCount: int64(len(result)),
	}, nil
}

func matchPart(part *inventory_v1.Part, filter *inventory_v1.PartsFilter) bool {
	if filter == nil {
		return true
	}

	// UUIDs
	if len(filter.Uuids) > 0 {
		found := false
		for _, u := range filter.Uuids {
			if part.Uuid == u {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Names
	if len(filter.Names) > 0 {
		found := false
		for _, n := range filter.Names {
			if part.Info.Name == n {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Categories
	if len(filter.Categories) > 0 {
		found := false
		for _, c := range filter.Categories {
			if part.Info.Category == c {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Manufacturer countries
	if len(filter.ManufacturerCountries) > 0 {
		found := false
		for _, country := range filter.ManufacturerCountries {
			if part.Info.Manufacturer.Country == country {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Tags
	if len(filter.Tags) > 0 {
		found := false
		for _, tag := range filter.Tags {
			for _, pt := range part.Info.Tags {
				if pt == tag {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func (s *InventoryService) DeletePart(ctx context.Context, req *inventory_v1.DeletePartRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	inventory, ok := s.inventory[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Part with uuid %s not found", req.GetUuid())
	}

	inventory.DeletedAt = timestamppb.New(time.Now())

	return &emptypb.Empty{}, nil
}

func (s *InventoryService) GetAll(ctx context.Context, req *inventory_v1.GetAllRequest) (*inventory_v1.GetAllResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []*inventory_v1.Part{}
	if len(s.inventory) == 0 {
		return nil, status.Error(codes.NotFound, "Inventories are empty")
	}
	for _, v := range s.inventory {
		result = append(result, v)
	}
	return &inventory_v1.GetAllResponse{
		Part:       result,
		TotalCount: int64(len(result)),
	}, nil
}

func (s *InventoryService) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	inventory, ok := s.inventory[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
	}

	return &inventory_v1.GetPartResponse{
		Part: inventory,
	}, nil
}

func (s *InventoryService) Create(ctx context.Context, req *inventory_v1.CreatePartsRequest) (*inventory_v1.CreatePartsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req == nil {
		log.Println("Received nil request")
		return nil, fmt.Errorf("request cannot be nil")
	}

	for _, v := range s.inventory {
		if v.Info.Name == req.GetInfo().Name {
			return nil, status.Error(codes.AlreadyExists, "part with this name already exists")
		}
	}

	newUUID := uuid.NewString()

	inventory := &inventory_v1.Part{
		Uuid:      newUUID,
		Info:      req.GetInfo(),
		CreatedAt: timestamppb.New(time.Now()),
	}

	s.inventory[newUUID] = inventory

	return &inventory_v1.CreatePartsResponse{
		Uuid: newUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
			// grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor2()),
		),
	)

	repo := inventory.NewRepository()
	service := service.NewInventoryServiceClient(repo)
	api := v1.NewInventoryApi(service)

	inventory_v1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	var gwServer *http.Server

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err = inventory_v1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx, mux, fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}
		gwServer = &http.Server{
			Addr:        fmt.Sprintf(":%d", httpPort),
			Handler:     mux,
			ReadTimeout: 10 * time.Second,
		}
		log.Printf("🌐 HTTP server with gRPC-Gateway listening on %d\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
