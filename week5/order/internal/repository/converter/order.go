package converter

import (
	"github.com/mbakhodurov/homeworks/week5/order/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week5/order/internal/repository/model"
)

func OrderUpdateInfoModelToRepoModel(from model.OrderUpdateInfo) repoModel.OrderUpdateInfo {
	return repoModel.OrderUpdateInfo{
		Status:           (*repoModel.OrderStatus)(from.Status),
		Transaction_uuid: from.Transaction_uuid,
		Payment_method:   (*repoModel.PaymentMethod)(from.Payment_method),
		Updated_at:       from.Updated_at,
		Deleted_at:       from.Deleted_at,
	}
}

func OrderRepoModelToModel(from repoModel.Order) model.Order {
	return model.Order{
		OrderUUID:       from.OrderUUID,
		UserUUID:        from.UserUUID,
		PartUUID:        from.PartUUID,
		TotalPrice:      from.TotalPrice,
		TransactionUUID: from.TransactionUUID,
		PaymentMethod:   (*model.PaymentMethod)(from.PaymentMethod),
		Status:          model.OrderStatus(from.Status),
		CreatedAt:       from.CreatedAt,
		UpdatedAt:       from.UpdatedAt,
		Deleted_at:      from.Deleted_at,
	}
}

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
		Deleted_at:      from.Deleted_at,
	}
}
