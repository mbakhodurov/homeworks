package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/mbakhodurov/homeworks/week2/order/internal/converter"
	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params order_v1.GetOrderByUUIDParams) (order_v1.GetOrderByUUIDRes, error) {
	order, err := a.OrderService.Get(ctx, params.OrderUUID)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Заказов пока что нету",
			}, nil
		}
		return &order_v1.InternalServerError{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &order_v1.GetOrderResponse{
		OrderDto: converter.OrderToDTO(order),
	}, nil

	// dto := order_v1.OrderDto{
	// 	OrderUUID:  order.Order_uuid,
	// 	UserUUID:   order.User_uuid,
	// 	PartUuids:  order.Part_uuids,
	// 	TotalPrice: float32(order.Total_price),
	// }
	// switch order.Status {
	// case model.StatusPendingPayment:
	// 	dto.Status = order_v1.OrderStatusPENDINGPAYMENT
	// case model.StatusPaid:
	// 	dto.Status = order_v1.OrderStatusPAID
	// case model.StatusCancelled:
	// 	dto.Status = order_v1.OrderStatusCANCELLED
	// default:
	// 	dto.Status = order_v1.OrderStatusPENDINGPAYMENT
	// }

	// if order.Transaction_uuid != nil {
	// 	dto.TransactionUUID = order_v1.NewOptNilString(*order.Transaction_uuid)
	// }

	// if order.Payment_method != nil {
	// 	var pm order_v1.OrderDtoPaymentMethod

	// 	switch *order.Payment_method {
	// 	case model.PaymentMethodCard:
	// 		pm = order_v1.OrderDtoPaymentMethodCARD
	// 	case model.PaymentMethodSBP:
	// 		pm = order_v1.OrderDtoPaymentMethodSBP
	// 	case model.PaymentMethodCreditCard:
	// 		pm = order_v1.OrderDtoPaymentMethodCREDITCARD
	// 	case model.PaymentMethodInvestorMoney:
	// 		pm = order_v1.OrderDtoPaymentMethodINVESTORMONEY
	// 	default:
	// 		pm = order_v1.OrderDtoPaymentMethodUNKNOWN
	// 	}

	// 	dto.PaymentMethod = &order_v1.NilOrderDtoPaymentMethod{
	// 		Value: pm,
	// 	}
	// }

	// return &order_v1.GetOrderResponse{
	// 	OrderDto: dto,
	// }, nil
}
