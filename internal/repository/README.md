# Package repository

This package is responsible for all database interactions.  Database functions should do the following:

- Use prepared statments when possible
- Limit business logic to a minimum
- Limit each function to one database interaction - if multiple interactions need to take place (e.g., A SQL INSERT and a subsequent SQL SELECT), those seperate functions can be called at the service layer
- Query functions should include a context.Context object

## Transactions

Transactions are managed and fully executed within this layer.  Transcations give us the ability to create/update multiple rows across any tables in our database and only commit the changes if there are no errors.  This will also give us loggable errors if we run into an issue and have to perform a rollback.

Since the entire transaction lifecycle takes place within the repository layer, that does introduce some opportunities for business logic to bleed into this layer.  Business logic should be kept to a minimum if possible.

> Note: We can also start a transaction and call repository functions completely in the service layer.  This is the preferred method since it doesn't create duplicate code and doesn't let business logic bleed into the repository layer.  Reference the [service layer docs](./../service/README.md) for more information

Transaction functions must include the following:
- Transaction function names must end in Tx
- Transaction functions must include a context.Context object
- Transaction Rollbacks must include an error message stating what failed

### Transaction Code Example 

```go
func (r *mysqlRepository) CreateCustomerTx(ctx context.Context, req model.Customer) (int64, error) {

    // Start the transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to begin transaction: %w", err)
	}
	defer func() {
        // This is where we perform transaction rollback.  
        // If we are here we ran into an error during the transaction

		slog.Warn("tranaction failed, rolling back: %w", err)
		tx.Rollback()
	}()

	query := "INSERT INTO customers (name) VALUES (?)"
	result, err := tx.ExecContext(ctx, query, req.Name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert data: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get insert ID: %w", err)
	}

	return id, nil
}
```