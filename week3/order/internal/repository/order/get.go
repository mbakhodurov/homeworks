package order

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	"github.com/mbakhodurov/homeworks/week3/order/internal/repository/converter"
	repoModel "github.com/mbakhodurov/homeworks/week3/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, order_uuid string) (model.Order, error) {
	builderSelect := squirrel.Select("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status", "created_at", "updated_at", "deleted_at").
		From("orders").PlaceholderFormat(squirrel.Dollar).Where(squirrel.Eq{"order_uuid": order_uuid}).OrderBy("id ASC")

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return model.Order{}, err
	}

	var order repoModel.Order

	err = r.db.
		QueryRowContext(ctx, query, args...).
		Scan(
			&order.OrderUUID,
			&order.UserUUID,
			pq.Array(&order.PartUUIDs),
			&order.TotalPrice,
			&order.TransactionUUID,
			&order.PaymentMethod,
			&order.Status,
			&order.CreatedAt,
			&order.Updated_at,
			&order.Deleted_at,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}

		log.Printf("failed to select order: %v\n", err)
		return model.Order{}, err
	}
	return converter.OrderRepoModelToModel(order), nil
}
