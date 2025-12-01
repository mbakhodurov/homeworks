package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	inventory_v1 "github.com/mbakhodurov/homeworks/week1/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
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

}
