package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"strconv"
)

func (s *Storage) AddCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount) error {
	query := "INSERT INTO accounts (Login, Password, Current, Withdraw) VALUES($1, $2, $3, $4)"
	balance := strconv.FormatFloat(customerAccount.Current, 'f', -1, 64)
	withdraw := strconv.FormatFloat(customerAccount.Withdraw, 'f', -1, 64)
	var args = []interface{}{customerAccount.Login, customerAccount.Password, balance, withdraw}
	return s.modifyCustomerAccount(ctx, customerAccount, query, args)
}

func (s *Storage) UpdateCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount) error {
	query := "UPDATE accounts SET Login = $1, Password = $2, Current = $3, Withdraw = $4 WHERE Login = $5"
	balance := strconv.FormatFloat(customerAccount.Current, 'f', -1, 64)
	withdraw := strconv.FormatFloat(customerAccount.Withdraw, 'f', -1, 64)
	var args = []interface{}{customerAccount.Login, customerAccount.Password, balance, withdraw, customerAccount.Login}
	return s.modifyCustomerAccount(ctx, customerAccount, query, args)
}

func (s *Storage) modifyCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount,
	query string, args []interface{}) error {
	logging.Debug("Preparing statement to update customer account to Repository: %+v", customerAccount)
	statement, err := s.Tx[ctx].Prepare(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return err
	}
	logging.Debug("Executing statement to update customer account to Repository: %+v", customerAccount)
	_, err = statement.Exec(args...)
	if err != nil {
		logging.Warn("Error During statement Execution %s with %s, %s, %s, %s",
			query, args[0], args[1], args[2], args[3])
		return err
	}
	return nil
}

func (s *Storage) GetCustomerAccount(login string) (*account.CustomerAccount, error) {
	return s.getCustomerAccountCommon(login, s.DB.Prepare)
}

func (s *Storage) GetCustomerAccountWithinTransaction(ctx context.Context, login string) (*account.CustomerAccount, error) {
	return s.getCustomerAccountCommon(login, s.Tx[ctx].Prepare)
}

func (s *Storage) getCustomerAccountCommon(login string, prepareFunc func(query string) (*sql.Stmt, error)) (*account.CustomerAccount, error) {
	query := "SELECT * FROM accounts WHERE Login = $1"
	var args = []interface{}{login}
	logging.Debug("Preparing statement to fetch customer account to Repository: %s", login)
	statement, err := prepareFunc(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return nil, err
	}
	logging.Debug("Executing statement to add customer account to Repository: %s", login)
	row := statement.QueryRow(args...)

	if row.Err() != nil {
		logging.Warn("Error During statement Execution %s with %s", query, login)
		return nil, err
	}
	a := &account.CustomerAccount{}
	err = row.Scan(&a.Login, &a.Password, &a.Current, &a.Withdraw)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		logging.Warn("Database response is empty")
		return nil, repositoryerrors.ErrNoContent
	case err != nil:
		logging.Warn("Error during decoding database response")
		return nil, err
	default:
		return a, nil
	}
}

func (s *Storage) GetCustomerAccounts() (account.CustomerAccounts, error) {
	return account.CustomerAccounts{}, nil
}

func (s *Storage) GetOrderOwner(ctx context.Context, orderID string) (string, error) {
	query := "SELECT owner FROM orders WHERE number = $1"
	var args = []interface{}{orderID}
	logging.Debug("Preparing statement to fetch customer account to Repository: %s", query)
	statement, err := s.Tx[ctx].Prepare(query)
	if err != nil {
		logging.Warn("Error During statement creation %s", query)
		return "", err
	}
	logging.Debug("Executing statement to add customer account to Repository: %s %s", query, args)
	row := statement.QueryRow(args...)

	if row.Err() != nil {
		logging.Warn("Error During statement Execution %s with %s", query, orderID)
		return "", err
	}
	var owner string
	err = row.Scan(&owner)
	if err != nil {
		logging.Warn("Could not decode Database Server response: %s", err.Error())
		return "", err
	}
	return owner, nil
}
