# Package service

This package is responisble for all business logic for our routes

## Transactions

Database calls can be wrapped in a transaction function in this layer using the `transactor.WithTransaction` function.  This allows us to keep business logic contained within the service layer but gives us the data integrity that comes with using a transaction.

### Transaction Code Example

```go
func (s *CustomerService) CreateCustomer(ctx context.Context, name string) (*model.Customer, error) {
	var customer *model.Customer

    // start the transaction
	err := s.transactor.WithTransaction(ctx, func(txRepo repository.Repository) error {
		_, err := txRepo.CreateCustomer(ctx, model.Customer{Name: name})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return customer, nil
}
```

> Note: WithTransaction only returns an error type (did the transaction suceed or not).  Since that is the case, you will need to think about what actions you are performing within the transaction and what types of data need to be returned.  If all you care about is if it completed or not, you might not have to do anything, but if you need to return an object you created/updated, you might need to make a follow-up database call.