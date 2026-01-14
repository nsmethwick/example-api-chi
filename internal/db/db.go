package db

import (
	"database/sql"
	"fmt"
	"time"
)

var DB *sql.DB

type Row struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InsertRequest struct {
	Name string `json:"name"`
}

type InsertResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func Init(dsn string) error {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		DB.Close() // best-effort cleanup
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}
	return DB.Close()
}

func GetAllCustomers() ([]Row, error) {
	query := "SELECT id, name, created_at, updated_at FROM customers"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var results []Row

	for rows.Next() {
		var id int
		var name string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &name, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, Row{
			ID:        id,
			Name:      name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return results, nil
}

func GetCustomerByID(id string) ([]Row, error) {
	query := "SELECT id, name, created_at, updated_at FROM customers WHERE id = ?"
	rows, err := DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var results []Row

	for rows.Next() {
		var rowID int
		var name string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&rowID, &name, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, Row{
			ID:        rowID,
			Name:      name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return results, nil
}

func CreateCustomer(req InsertRequest) (InsertResponse, error) {
	name := req.Name
	if name == "" {
		return InsertResponse{}, fmt.Errorf("Name field missing")
	}

	query := "INSERT INTO customers (name) VALUES (?)"
	result, err := DB.Exec(query, name)
	if err != nil {
		return InsertResponse{}, fmt.Errorf("failed to insert data: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return InsertResponse{}, fmt.Errorf("failed to get insert ID: %w", err)
	}

	return InsertResponse{
		ID:   id,
		Name: name,
	}, nil
}

// Adding this as an option
func CreateCustomerFullResponse(req InsertRequest) (Row, error) {
	name := req.Name
	if name == "" {
		return Row{}, fmt.Errorf("Name field missing")
	}

	query := "INSERT INTO customers (name) VALUES (?)"
	result, err := DB.Exec(query, name)
	if err != nil {
		return Row{}, fmt.Errorf("failed to insert data: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Row{}, fmt.Errorf("failed to get insert ID: %w", err)
	}

	var r Row
	q := "SELECT id, name, created_at, updated_at FROM customers WHERE id = ?"
	if err := DB.QueryRow(q, id).Scan(&r.ID, &r.Name, &r.CreatedAt, &r.UpdatedAt); err != nil {
		return Row{}, fmt.Errorf("failed to fetch inserted row: %w", err)
	}

	return r, nil
}
