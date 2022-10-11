package database

import "github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"

func (s *Storage) AddCustomerAccount(customerAccount *customer_account.CustomerAccount) error {
	return nil
}

func (s *Storage) UpdateCustomerAccount(customerAccount *customer_account.CustomerAccount) error {
	return nil
}

func (s *Storage) GetCustomerAccount(login string) (*customer_account.CustomerAccount, error) {
	return nil, nil
}

func (s *Storage) GetCustomerAccounts() (customer_account.CustomerAccounts, error) {
	return customer_account.CustomerAccounts{}, nil
}

func (s *Storage) GetOrderOwner(orderId string) (string, error) {
	return "", nil
}
