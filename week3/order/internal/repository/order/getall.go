package order

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	"github.com/mbakhodurov/homeworks/week3/order/internal/repository/converter"
	repoModel "github.com/mbakhodurov/homeworks/week3/order/internal/repository/model"
)

func (r *repository) GetAll(ctx context.Context) ([]model.Order, error) {
	builderSelect := squirrel.Select("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status", "created_at", "updated_at", "deleted_at").
		From("orders").PlaceholderFormat(squirrel.Dollar).
		OrderBy("id ASC")

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return []model.Order{}, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("failed to select orders: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var o repoModel.Order

		err := rows.Scan(
			&o.OrderUUID,
			&o.UserUUID,
			pq.Array(&o.PartUUIDs),
			&o.TotalPrice,
			&o.TransactionUUID,
			&o.PaymentMethod,
			&o.Status,
			&o.CreatedAt,
			&o.Updated_at,
			&o.Deleted_at,
		)
		if err != nil {
			log.Printf("failed to scan order: %v\n", err)
			return nil, err
		}

		converted := converter.OrderRepoModelToModel(o)

		orders = append(orders, converted)

	}

	if len(orders) == 0 {
		return nil, model.ErrOrderNotFound
	}

	return orders, nil
}
