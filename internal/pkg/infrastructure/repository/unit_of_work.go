package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	db "golang-boilerplate/internal/pkg/connections/db"
)

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(uow UnitOfWork) error) (err error)
	ProductRepo() ProductRepository
	BillerRepo() BillerRepository
	ProductBillerRepo() ProductBillerRepository
}

type unitOfWork struct {
	db db.DBExecutor
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

func (uow *unitOfWork) Execute(ctx context.Context, fn func(uow UnitOfWork) error) error {
	return db.WithTransaction(ctx, uow.db, func(tx *sqlx.Tx) error {
		return fn(&unitOfWork{
			db: tx,
		})
	})
}

func (uow *unitOfWork) ProductRepo() ProductRepository {
	return NewProductRepository(uow.db)
}

func (uow *unitOfWork) BillerRepo() BillerRepository {
	return NewBillerRepository(uow.db)
}

func (uow *unitOfWork) ProductBillerRepo() ProductBillerRepository {
	return NewProductBillerRepository(uow.db)
}
