package memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
)

func (s *Storage) GetCustomerAccount(login string) (*account.CustomerAccount, error) {
	res, ok := s.CustomerAccounts[login]
	if !ok {
		return nil, repositoryerrors.ErrNoContent
	}
	return &res, nil
}

func (s *Storage) GetCustomerAccountWithinTransaction(ctx context.Context, login string) (*account.CustomerAccount, error) {
	return s.GetCustomerAccount(login)
}

func (s *Storage) GetCustomerAccounts() (account.CustomerAccounts, error) {
	return s.CustomerAccounts, nil
}

func (s *Storage) AddCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount) error {
	_, ok := s.CustomerAccounts[customerAccount.Login]
	if !ok {
		s.CustomerAccounts[customerAccount.Login] = *customerAccount
		return nil
	}
	return fmt.Errorf("record Alreadt present")
}

func (s *Storage) UpdateCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount) error {
	_, exists := s.CustomerAccounts[customerAccount.Login]
	if !exists {
		return errors.New("relevant record does not exist")
	}
	s.CustomerAccounts[customerAccount.Login] = *customerAccount
	return nil
}

func (s *Storage) GetOrderOwner(ctx context.Context, orderID string) (string, error) {
	order, exists := s.Orders[orderID]
	if !exists {
		return "", fmt.Errorf("record does not exist")
	}
	return order.Owner, nil
}
