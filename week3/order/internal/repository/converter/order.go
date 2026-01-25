package converter

import (
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	repo "github.com/mbakhodurov/homeworks/week3/order/internal/repository/model"
)

func OrderUpdateInfoModelToRepoModel(from model.OrderUpdateInfo) repo.OrderUpdateInfo {
	return repo.OrderUpdateInfo{
		Status:           (*repo.OrderStatus)(from.Status),
		Transaction_uuid: from.Transaction_uuid,
		Payment_method:   (*repo.PaymentMethod)(from.Payment_method),
		Updated_at:       from.Updated_at,
		Deleted_at:       from.Deleted_at,
	}
}

func OrderRepoModelToModel(from repo.Order) model.Order {
	return model.Order{
		OrderUUID:       from.OrderUUID,
		UserUUID:        from.UserUUID,
		PartUUIDs:       from.PartUUIDs,
		TotalPrice:      from.TotalPrice,
		TransactionUUID: from.TransactionUUID,
		PaymentMethod:   (*model.PaymentMethod)(from.PaymentMethod),
		Status:          model.OrderStatus(from.Status),
		CreatedAt:       from.CreatedAt,
		Updated_at:      from.Updated_at,
		Deleted_at:      from.Deleted_at,
	}
}

func OrderModelToRepoModel(from model.Order) repo.Order {
	return repo.Order{
		OrderUUID:       from.OrderUUID,
		UserUUID:        from.UserUUID,
		PartUUIDs:       from.PartUUIDs,
		TotalPrice:      from.TotalPrice,
		TransactionUUID: from.TransactionUUID,
		PaymentMethod:   (*repo.PaymentMethod)(from.PaymentMethod),
		Status:          repo.OrderStatus(from.Status),
		CreatedAt:       from.CreatedAt,
		Updated_at:      from.Updated_at,
		Deleted_at:      from.Deleted_at,
	}
}
