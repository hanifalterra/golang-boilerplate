package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"golang-boilerplate/internal/pkg/connections/database"
)

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(uow UnitOfWork) error) (err error)
	GetProductBillerRepo() ProductBillerRepository
}

type unitOfWork struct {
	db database.DBExecutor
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

func (uow *unitOfWork) Execute(ctx context.Context, fn func(uow UnitOfWork) error) error {
	return database.WithTransaction(ctx, uow.db, "UnitOfWork", "Execute", func(tx *sqlx.Tx) error {
		return fn(&unitOfWork{
			db: tx,
		})
	})
}

func (uow *unitOfWork) GetProductBillerRepo() ProductBillerRepository {
	return NewProductBillerRepository(uow.db)
}
