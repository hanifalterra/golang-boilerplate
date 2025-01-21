package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/models"
)

func TestProductBillerRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlx.NameMapper = strcase.ToSnake
	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repositories.NewProductBillerRepository(sqlxDB)

	productBiller := &models.ProductBiller{
		ProductID: 1,
		BillerID:  2,
		IsActive:  true,
		CreatedBy: "test_user",
		UpdatedBy: "test_user",
	}

	query := `
		INSERT INTO product_billers 
		\(product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by\)
		VALUES \(\?, \?, \?, NOW\(6\), \?, NOW\(6\), \?\)
	`
	mock.ExpectExec(query).
		WithArgs(
			productBiller.ProductID,
			productBiller.BillerID,
			productBiller.IsActive,
			productBiller.CreatedBy,
			productBiller.UpdatedBy,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), productBiller)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductBillerRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repositories.NewProductBillerRepository(sqlxDB)

	query := `UPDATE product_billers SET is_active = \?, updated_at = NOW\(6\) WHERE id = \? AND deleted_at IS NULL`
	mock.ExpectExec(query).
		WithArgs(false, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(context.Background(), 1, &models.ProductBiller{IsActive: false})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductBillerRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repositories.NewProductBillerRepository(sqlxDB)

	query := `UPDATE product_billers SET deleted_at = NOW\(6\) WHERE id = \? AND deleted_at IS NULL`
	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductBillerRepository_FetchOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlx.NameMapper = strcase.ToSnake
	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repositories.NewProductBillerRepository(sqlxDB)

	query := `SELECT id, product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by FROM product_billers WHERE deleted_at IS NULL AND id = \?`
	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "biller_id", "is_active", "created_at", "created_by", "updated_at", "updated_by"}).
			AddRow(1, 1, 2, true, time.Time{}, "user1", time.Time{}, "user1"))

	result, err := repo.FetchOne(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductBillerRepository_FetchMany(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlx.NameMapper = strcase.ToSnake
	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repositories.NewProductBillerRepository(sqlxDB)

	query := `SELECT id, product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by FROM product_billers WHERE deleted_at IS NULL AND product_id = \?`
	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "biller_id", "is_active", "created_at", "created_by", "updated_at", "updated_by"}).
			AddRow(1, 1, 2, true, time.Time{}, "user1", time.Time{}, "user1").
			AddRow(2, 1, 3, false, time.Time{}, "user2", time.Time{}, "user2"))

	results, err := repo.FetchMany(context.Background(), map[string]interface{}{"product_id": 1})
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
