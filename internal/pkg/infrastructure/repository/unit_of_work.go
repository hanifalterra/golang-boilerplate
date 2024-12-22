package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	db "golang-boilerplate/internal/pkg/connections/db"
)

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(uow UnitOfWork) error) (err error)
	GetProductBillerRepo() ProductBillerRepository
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
	return db.WithTransaction(ctx, uow.db, "UnitOfWork", "Execute", func(tx *sqlx.Tx) error {
		return fn(&unitOfWork{
			db: tx,
		})
	})
}

func (uow *unitOfWork) GetProductBillerRepo() ProductBillerRepository {
	return NewProductBillerRepository(uow.db)
}
