package v1

import inventory_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/proto/inventory/v1"

type Client struct {
	generatedClient inventory_v1.InventoryServiceClient
}

func NewClient(generatedClient inventory_v1.InventoryServiceClient) *Client {
	return &Client{
		generatedClient: generatedClient,
	}
}
