package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/pick-cee/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

var (
	InvalidOrderError = errors.New("Invalid order data.")
	ProductNotFoundError = errors.New("Product not found")
	ProductNoStockError = errors.New("No available stock for this product")
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderRequest) (repo.Order, error)
}

type svc struct {
	repo *repo.Queries
	db *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db: db,
	}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderRequest) (repo.Order, error) {
	if tempOrder.CustomerID == 0 || len(tempOrder.Items) == 0 {
		return repo.Order{}, InvalidOrderError
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create an order
	order, err := qtx.CreateOrder(ctx, int64(tempOrder.CustomerID))
	if err != nil {
		return repo.Order{}, err
	}

	// Look for product if exists 
	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductById(ctx, int64(item.ProductID))
		if err != nil {
			return repo.Order{}, ProductNotFoundError
		}

		if product.Quantity < int32(item.Quantity) {
			return repo.Order{}, ProductNoStockError
		}

		// create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID: order.ID,
			ProductID: int64(item.ProductID),
			Quantity: int32(item.Quantity),
			PriceCents: product.PriceInCents,
		})

		if err != nil {
			return repo.Order{}, err
		}

		// update the product quantity
		qtx.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
			ID: int64(item.ProductID),
			Quantity: product.Quantity - int32(item.Quantity),
		})
	}

	tx.Commit(ctx)
	return order, nil
}