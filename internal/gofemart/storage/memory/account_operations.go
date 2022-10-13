package memory

import (
	"errors"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repository_errors"
)

func (s *Storage) GetCustomerAccount(login string) (*customer_account.CustomerAccount, error) {
	res, ok := s.CustomerAccounts[login]
	if !ok {
		return nil, repository_errors.ErrNoContent
	}
	return &res, nil
}

func (s *Storage) GetCustomerAccountWithinTransaction(login string) (*customer_account.CustomerAccount, error) {
	return s.GetCustomerAccount(login)
}

func (s *Storage) GetCustomerAccounts() (customer_account.CustomerAccounts, error) {
	return s.CustomerAccounts, nil
}

func (s *Storage) AddCustomerAccount(customerAccount *customer_account.CustomerAccount) error {
	_, ok := s.CustomerAccounts[customerAccount.Login]
	if !ok {
		s.CustomerAccounts[customerAccount.Login] = *customerAccount
		return nil
	}
	return fmt.Errorf("record Alreadt present")
}

func (s *Storage) UpdateCustomerAccount(customerAccount *customer_account.CustomerAccount) error {
	_, exists := s.CustomerAccounts[customerAccount.Login]
	if !exists {
		return errors.New("relevant record does not exist")
	}
	s.CustomerAccounts[customerAccount.Login] = *customerAccount
	return nil
}

func (s *Storage) GetOrderOwner(orderId string) (string, error) {
	order, exists := s.Orders[orderId]
	if !exists {
		return "", fmt.Errorf("record does not exist")
	}
	return order.Owner, nil
}
