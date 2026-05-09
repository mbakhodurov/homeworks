package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/order/internal/model"
)

func (s *service) GetAll(ctx context.Context) ([]model.Order, error) {
	// sessionUUID, ok := httpMiddleware.GetUserFromContext(ctx)
	// if ok {
	// 	fmt.Println("SESSION:", sessionUUID)
	// }
	res, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return []model.Order{}, err
	}

	return res, nil
}
