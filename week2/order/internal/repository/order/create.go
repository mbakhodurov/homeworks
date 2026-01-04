package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week2/order/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[order.Order_uuid] = repoModel.Order{
		Order_uuid:       order.Order_uuid,
		User_uuid:        order.User_uuid,
		Part_uuids:       order.Part_uuids,
		Total_price:      order.Total_price,
		Transaction_uuid: order.Transaction_uuid,
		Payment_method:   (*repoModel.PaymentMethod)(order.Payment_method),
		Status:           repoModel.OrderStatus(order.Status),
		Created_at:       order.Created_at,
	}
	return nil
}
