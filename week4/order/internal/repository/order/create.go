package order

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/mbakhodurov/homeworks/week4/order/internal/model"
	"github.com/mbakhodurov/homeworks/week4/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) (int64, error) {
	repoModel := converter.OrderModelToRepoModel(order)
	builderInsert := squirrel.Insert("orders").
		PlaceholderFormat(squirrel.Dollar).
		Columns("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status", "created_at", "updated_at").
		Values(repoModel.OrderUUID, repoModel.UserUUID, pq.Array(repoModel.PartUUID), repoModel.TotalPrice, repoModel.TransactionUUID, repoModel.PaymentMethod, repoModel.Status, repoModel.CreatedAt, repoModel.UpdatedAt).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return 0, err
	}

	var orderId int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&orderId)
	if err != nil {
		log.Printf("failed to insert note: %v\n", err)
		return 0, err
	}
	return int64(orderId), nil
}
