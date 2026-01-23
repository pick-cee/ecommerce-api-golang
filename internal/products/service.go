package products

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/pick-cee/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

var (
	InvalidProductRequest = errors.New("Invalid Product request data, please recheck.")
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	CreateProduct(ctx context.Context, tempProduct createProduct) (repo.Product, error)
	UpdateProductQuantity(ctx context.Context, productId int, newQuantity int) (repo.Product, error)
	DeleteProduct(ctx context.Context, productId int) error
}

type svc struct {
	repo *repo.Queries
	db *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{repo: repo, db: db}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	product, err := s.repo.ListProducts(ctx)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *svc) CreateProduct(ctx context.Context, tempProduct createProduct) (repo.Product, error) {
	if tempProduct.Name == "" || tempProduct.PriceInCents == 0 || tempProduct.Quantity == 0 {
		return repo.Product{}, InvalidProductRequest
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Product{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create the product
	product, err := qtx.CreateProduct(ctx, repo.CreateProductParams{
		Name: tempProduct.Name,
		Quantity: int32(tempProduct.Quantity),
		PriceInCents: int32(tempProduct.PriceInCents),
	})

	if err != nil {
		return repo.Product{}, err
	}

	tx.Commit(ctx)
	
	return  product, nil
}

func (s *svc) UpdateProductQuantity(ctx context.Context, productId int, newQuantity int) (repo.Product, error) {
	currentProduct, err := s.repo.FindProductById(ctx, int64(productId))

	if err != nil {
		return repo.Product{}, err
	}

	updatedProduct, err := s.repo.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
		ID: int64(productId),
		Quantity: int32(newQuantity) + currentProduct.Quantity,
	})

	if err != nil {
		return repo.Product{}, nil
	}

	return  updatedProduct, nil
}

func (s *svc) DeleteProduct(ctx context.Context, productId int) error {
	_, err := s.repo.FindProductById(ctx, int64(productId))
	if err != nil {
		return err
	}

	return s.repo.DeleteProduct(ctx, int64(productId))
}