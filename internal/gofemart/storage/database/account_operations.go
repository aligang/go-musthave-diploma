package database

import (
	"context"
	"database/sql"

	"errors"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/repositoryerrors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

func (s *Storage) AddCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount, tx *sqlx.Tx) error {
	query := "INSERT INTO accounts (Login, Password, Current, Withdraw) VALUES($1, $2, $3, $4)"
	balance := strconv.FormatFloat(customerAccount.Current, 'f', -1, 64)
	withdraw := strconv.FormatFloat(customerAccount.Withdraw, 'f', -1, 64)
	var args = []interface{}{customerAccount.Login, customerAccount.Password, balance, withdraw}
	return s.modifyCustomerAccount(ctx, customerAccount, query, args, tx)
}

func (s *Storage) UpdateCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount, tx *sqlx.Tx) error {
	query := "UPDATE accounts SET Login = $1, Password = $2, Current = $3, Withdraw = $4 WHERE Login = $5"
	balance := strconv.FormatFloat(customerAccount.Current, 'f', -1, 64)
	withdraw := strconv.FormatFloat(customerAccount.Withdraw, 'f', -1, 64)
	var args = []interface{}{customerAccount.Login, customerAccount.Password, balance, withdraw, customerAccount.Login}
	return s.modifyCustomerAccount(ctx, customerAccount, query, args, tx)
}

func (s *Storage) modifyCustomerAccount(ctx context.Context, customerAccount *account.CustomerAccount,
	query string, args []interface{}, tx *sqlx.Tx) error {
	s.log.Debug("Preparing statement to update customer account to Repository: %+v", customerAccount)
	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return err
	}
	s.log.Debug("Executing statement to update customer account to Repository: %+v", customerAccount)
	_, err = statement.ExecContext(ctx, args...)
	if err != nil {
		s.log.Warn("Error During statement Execution %s with %s, %s, %s, %s",
			query, args[0], args[1], args[2], args[3])
		return err
	}
	return nil
}

func (s *Storage) GetCustomerAccount(ctx context.Context, login string, tx *sqlx.Tx) (*account.CustomerAccount, error) {
	query := "SELECT * FROM accounts WHERE Login = $1"
	var args = []interface{}{login}
	s.log.Debug("Preparing statement to fetch customer account to Repository: %s", login)
	var err error
	var statement *sqlx.Stmt
	a := &account.CustomerAccount{}

	if tx != nil {
		statement, err = tx.PreparexContext(ctx, query)
	} else {
		statement, err = s.DB.PreparexContext(ctx, query)
	}
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return nil, err
	}
	s.log.Debug("Executing statement to add customer account to Repository: %s", login)
	err = statement.GetContext(ctx, a, args...)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		s.log.Warn("Database response is empty")
		return nil, repositoryerrors.ErrNoContent
	case err != nil:
		s.log.Warn("Error during decoding database response")
		return nil, err
	default:
		return a, nil
	}
}

func (s *Storage) GetCustomerAccounts() (account.CustomerAccounts, error) {
	return account.CustomerAccounts{}, nil
}

func (s *Storage) GetOrderOwner(ctx context.Context, orderID string, tx *sqlx.Tx) (string, error) {
	query := "SELECT owner FROM orders WHERE number = $1"
	var args = []interface{}{orderID}
	s.log.Debug("Preparing statement to fetch customer account to Repository: %s", query)
	statement, err := tx.PreparexContext(ctx, query)
	if err != nil {
		s.log.Warn("Error During statement creation %s", query)
		return "", err
	}
	s.log.Debug("Executing statement to add customer account to Repository: %s %s", query, args)
	var owner string
	err = statement.GetContext(ctx, &owner, args...)

	if err != nil {
		s.log.Warn("Could not decode Database Server response: %s", err.Error())
		return "", err
	}
	return owner, nil
}
