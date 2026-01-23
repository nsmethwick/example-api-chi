package repository

import (
	"context"
	"database/sql"
	"example-api-chi/internal/model"
	"fmt"
	"time"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// TODO:
// Consider renaming these, we will hold all db interactions here and want to make sure the names don't collide in the future
type Repository interface {
	GetCustomers(ctx context.Context) ([]model.Customer, error)
	GetCustomerByID(ctx context.Context, id string) (model.Customer, error)
	CreateCustomer(ctx context.Context, req model.Customer) (int64, error)
	UpdateCustomer(ctx context.Context, req model.Customer, id string) (int, error)
	DeleteCustomer(ctx context.Context, id string) error
}

type mysqlRepository struct {
	db DBTX
}

func NewRepository(database DBTX) (Repository, error) {
	return &mysqlRepository{
		db: database,
	}, nil
}

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) (*Transactor, error) {
	return &Transactor{db: db}, nil
}

func (t *Transactor) WithTransaction(ctx context.Context, fn func(repo Repository) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	txRepo, err := NewRepository(tx)
	if err := fn(txRepo); err != nil {
		fmt.Print(err)
		return err
	}
	return tx.Commit()
}

func (r *mysqlRepository) GetCustomers(ctx context.Context) ([]model.Customer, error) {
	query := "SELECT id, name, created_at, updated_at FROM customers"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var results []model.Customer

	for rows.Next() {
		var id int
		var name string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &name, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, model.Customer{
			ID:        id,
			Name:      name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return results, nil
}

func (r *mysqlRepository) GetCustomerByID(ctx context.Context, id string) (model.Customer, error) {
	var results model.Customer
	query := "SELECT id, name, created_at, updated_at FROM customers WHERE id = ?"
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&results.ID, &results.Name, &results.CreatedAt, &results.UpdatedAt); err != nil {
		return model.Customer{}, fmt.Errorf("failed to fetch inserted row: %w", err)
	}

	return results, nil
}

func (r *mysqlRepository) CreateCustomer(ctx context.Context, req model.Customer) (int64, error) {
	query := "INSERT INTO customers (name) VALUES (?)"
	result, err := r.db.ExecContext(ctx, query, req.Name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert data: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get insert ID: %w", err)
	}

	return id, nil
}

func (r *mysqlRepository) UpdateCustomer(ctx context.Context, req model.Customer, id string) (int, error) {
	query := "UPDATE customers SET name = ? WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, req.Name, id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert data: %w", err)
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get insert ID: %w", err)
	}

	return int(rowCount), nil
}

func (r *mysqlRepository) DeleteCustomer(ctx context.Context, id string) error {
	query := "DELETE FROM customers WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}

	return nil
}
