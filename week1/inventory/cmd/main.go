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
	inventory_v1 "github.com/mbakhodurov/homeworks/week1/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort = 50052
	httpPort = 8082
)

type InventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer
	mu sync.RWMutex

	inventories map[string]*inventory_v1.Part
}

func (i *InventoryService) GetPart(ctx context.Context, rq *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	part, ok := i.inventories[rq.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", rq.GetUuid())
	}

	return &inventory_v1.GetPartResponse{
		Part: part,
	}, nil
}

func (i *InventoryService) CreateParts(ctx context.Context, rq *inventory_v1.CreatePartsRequest) (*inventory_v1.CreatePartsResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if rq == nil {
		log.Println("Received nil request")
		return nil, fmt.Errorf("request cannot be nil")
	}

	for _, v := range i.inventories {
		if v.Info.Name == rq.GetInfo().Name {
			return nil, status.Error(codes.AlreadyExists, "part with this name already exists")
		}
	}

	newUUID := uuid.NewString()

	inventory := &inventory_v1.Part{
		Uuid:      newUUID,
		Info:      rq.GetInfo(),
		CreatedAt: timestamppb.New(time.Now()),
	}

	i.inventories[newUUID] = inventory

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
	s := grpc.NewServer()
	service := &InventoryService{
		inventories: make(map[string]*inventory_v1.Part),
	}
	inventory_v1.RegisterInventoryServiceServer(s, service)

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
