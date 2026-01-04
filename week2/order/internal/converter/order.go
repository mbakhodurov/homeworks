package converter

import (
	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/openapi/order/v1"
)

const (
	Unknown      = "unknown"
	Card         = "card"
	Sbp          = "sbp"
	CreditCard   = "credit-card"
	InvestorCard = "investor-money"
)

func PaymentMethodToModel(method order_v1.PaymentMethod) model.PaymentMethod {
	switch method {
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

func OrderStatusToAPI(status model.OrderStatus) order_v1.OrderStatus {
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

func OrderToDTO(order model.Order) order_v1.OrderDto {
	dto := order_v1.OrderDto{
		OrderUUID:  order.Order_uuid,
		UserUUID:   order.User_uuid,
		PartUuids:  order.Part_uuids,
		TotalPrice: float32(order.Total_price),
		Status:     OrderStatusToAPI(order.Status),
		CreatedAt:  order_v1.NewOptDateTime(order.Created_at),
	}

	if order.Transaction_uuid != nil {
		dto.TransactionUUID = order_v1.NewOptNilString(*order.Transaction_uuid)
	}

	if order.Payment_method != nil {
		apiPaymentMethod := PaymentMethodToOrderDtoPaymentMethod(*order.Payment_method)
		dto.PaymentMethod = &order_v1.NilOrderDtoPaymentMethod{
			Value: apiPaymentMethod,
			Null:  false,
		}
	}

	if !order.Deleted_at.IsZero() {
		dto.DeletedAt = order_v1.NewOptNilDateTime(order.Deleted_at)
	}

	if !order.Updated_at.IsZero() {
		dto.UpdatedAt = order_v1.NewOptNilDateTime(order.Updated_at)
	}

	return dto
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
