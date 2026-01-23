package service

import (
	"context"
	"example-api-chi/internal/model"
	"example-api-chi/internal/repository"
	"strconv"
)

type CustomerService struct {
	repo       repository.Repository
	transactor *repository.Transactor
}

func NewCustomerService(repo repository.Repository, transactor *repository.Transactor) *CustomerService {
	return &CustomerService{
		repo:       repo,
		transactor: transactor,
	}
}

func (s *CustomerService) GetAllCustomers(ctx context.Context) ([]model.Customer, error) {
	return s.repo.GetCustomers(ctx)
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, id string) (*model.Customer, error) {
	customers, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &customers, nil
}

func (s *CustomerService) CreateCustomer(ctx context.Context, name string) (*model.Customer, error) {
	id, err := s.repo.CreateCustomer(ctx, model.Customer{Name: name})
	if err != nil {
		return nil, err
	}

	stringValue := strconv.FormatInt(id, 10)
	customer, err := s.repo.GetCustomerByID(ctx, stringValue)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, id string, name string) (*model.Customer, error) {
	_, err := s.repo.UpdateCustomer(ctx, model.Customer{Name: name}, id)
	if err != nil {
		return nil, err
	}
	customer, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string) error {
	return s.repo.DeleteCustomer(ctx, id)
}
