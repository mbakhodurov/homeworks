package v1

import payment_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/payment/v1"

type client struct {
	generatedClient payment_v1.PaymentServiceClient
}

func NewClient(generatedClient payment_v1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
