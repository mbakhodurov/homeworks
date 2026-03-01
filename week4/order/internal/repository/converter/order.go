package converter

import (
	"github.com/mbakhodurov/homeworks/week4/order/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week4/order/internal/repository/model"
)

func OrderModelToRepoModel(from model.Order) repoModel.Order {
	return repoModel.Order{
		OrderUUID:       from.OrderUUID,
		UserUUID:        from.OrderUUID,
		PartUUID:        from.PartUUID,
		TotalPrice:      from.TotalPrice,
		TransactionUUID: from.TransactionUUID,
		PaymentMethod:   (*repoModel.PaymentMethod)(from.PaymentMethod),
		Status:          repoModel.OrderStatus(from.Status),
		CreatedAt:       from.CreatedAt,
		UpdatedAt:       from.UpdatedAt,
	}
}
