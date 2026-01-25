package converter

import (
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/openapi/order/v1"
)

func OrderModelToOrderDTO(from model.Order) order_v1.OrderDto {
	dto := order_v1.OrderDto{
		OrderUUID:  from.OrderUUID,
		UserUUID:   from.UserUUID,
		PartUuids:  from.PartUUIDs,
		TotalPrice: float32(from.TotalPrice),
		Status:     OrderStatusModelToDTO(from.Status),
	}

	if from.CreatedAt != nil {
		dto.CreatedAt = order_v1.NewOptDateTime(*from.CreatedAt)
	}

	if from.TransactionUUID != nil {
		dto.TransactionUUID = order_v1.NewOptNilString(*from.TransactionUUID)
	}

	if from.PaymentMethod != nil {
		apiPaymentMethod := PaymentMethodToOrderDtoPaymentMethod(*from.PaymentMethod)

		dto.PaymentMethod = &order_v1.NilOrderDtoPaymentMethod{
			Value: apiPaymentMethod,
			Null:  false,
		}
	}

	if from.Deleted_at != nil {
		dto.DeletedAt = order_v1.NewOptNilDateTime(*from.Deleted_at)
	}

	if from.Updated_at != nil {
		dto.UpdatedAt = order_v1.NewOptNilDateTime(*from.Updated_at)
	}

	return dto
}

func OrderStatusModelToDTO(status model.OrderStatus) order_v1.OrderStatus {
	switch status {
	case model.StatusPendingPayment:
		return order_v1.OrderStatusPENDINGPAYMENT
	case model.StatusPaid:
		return order_v1.OrderStatusPAID
	case model.StatusCancelled:
		return order_v1.OrderStatusCANCELLED
	case model.StatusCompleted:
		return order_v1.OrderStatusCOMPLETED
	default:
		return order_v1.OrderStatusUNKNOWN
	}
}

func PaymentMethodProtoToPaymentModel(protoMethod order_v1.PaymentMethod) model.PaymentMethod {
	switch protoMethod {
	case order_v1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case order_v1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case order_v1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case order_v1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func PaymentMethodToOrderDtoPaymentMethod(domainMethod model.PaymentMethod) order_v1.OrderDtoPaymentMethod {
	switch domainMethod {
	case model.PaymentMethodCard:
		return order_v1.OrderDtoPaymentMethodCARD
	case model.PaymentMethodSBP:
		return order_v1.OrderDtoPaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return order_v1.OrderDtoPaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return order_v1.OrderDtoPaymentMethodINVESTORMONEY
	default:
		return order_v1.OrderDtoPaymentMethodUNKNOWN
	}
}
