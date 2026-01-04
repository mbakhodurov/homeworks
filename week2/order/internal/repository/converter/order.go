package converter

import (
	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	repo "github.com/mbakhodurov/homeworks/week2/order/internal/repository/model"
)

// func OrderUpdateInfoToModel(order model.OrderUpdateInfo) repo.OrderUpdateInfo {
// 	return repo.OrderUpdateInfo{
// 		Payment_method: (*repo.PaymentMethod)(order.Payment_method),
// 		Updated_at:     order.Updated_at,
// 		Deleted_at:     order.Deleted_at,
// 		Status:         (*repo.OrderStatus)(order.Status),
// 	}
// }

func OrderToModel(order repo.Order) model.Order {
	return model.Order{
		Order_uuid:       order.Order_uuid,
		User_uuid:        order.User_uuid,
		Part_uuids:       order.Part_uuids,
		Total_price:      order.Total_price,
		Transaction_uuid: order.Transaction_uuid,
		Payment_method:   (*model.PaymentMethod)(order.Payment_method),
		Status:           model.OrderStatus(order.Status),
		Created_at:       order.Created_at,
		Updated_at:       order.Updated_at,
		Deleted_at:       order.Deleted_at,
	}
}
