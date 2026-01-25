package order

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	"github.com/mbakhodurov/homeworks/week3/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) (int64, error) {
	repoModel := converter.OrderModelToRepoModel(order)

	builderInsert := squirrel.Insert("orders").
		PlaceholderFormat(squirrel.Dollar).
		Columns("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status", "created_at", "updated_at").
		Values(repoModel.OrderUUID, repoModel.UserUUID, pq.Array(repoModel.PartUUIDs), repoModel.TotalPrice, repoModel.TransactionUUID, repoModel.PaymentMethod, repoModel.Status, repoModel.CreatedAt, repoModel.Updated_at).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return 0, err
	}

	var id int64
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		log.Printf("failed to execute query: %v\n", err)
		return 0, err
	}
	return id, nil
}
