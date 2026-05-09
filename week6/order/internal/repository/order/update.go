package order

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/mbakhodurov/homeworks/week6/order/internal/model"
	"github.com/mbakhodurov/homeworks/week6/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order_uuid string, updateInfo model.OrderUpdateInfo) error {
	repoModel := converter.OrderUpdateInfoModelToRepoModel(updateInfo)

	builderUpdate := squirrel.Update("orders").
		Where(squirrel.Eq{"order_uuid": order_uuid}).
		PlaceholderFormat(squirrel.Dollar).
		Set("status", repoModel.Status).
		Set("transaction_uuid", repoModel.Transaction_uuid).
		Set("payment_method", repoModel.Payment_method).
		Set("updated_at", repoModel.Updated_at).
		Set("Deleted_at", repoModel.Deleted_at)

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build update query: %v", err)
		return err
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update order: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrOrderNotFound
	}

	log.Printf("updated %d rows for order: %s", rowsAffected, order_uuid)
	return nil
}
