package database

import (
	"github.com/jmoiron/sqlx"

	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/logging"
)
import _ "github.com/jackc/pgx/v4/stdlib"

type Storage struct {
	DB  *sqlx.DB
	log *logging.InternalLogger
}

func New(conf *config.Config) *Storage {
	logging.Debug("Initialisating SQL Repository")
	db, err := sqlx.Open("pgx", conf.DatabaseURI)
	if err != nil {
		panic(err)
	}
	s := &Storage{
		DB:  db,
		log: logging.Logger.GetSubLogger("repository", "sql"),
	}
	_, err = s.DB.Exec(
		"create table if not exists accounts(Login text NOT NULL UNIQUE, Password text NOT NULL, Current double precision, Withdraw double precision)",
	)
	if err != nil {
		msg, _ := fmt.Printf("Failure during initialisation of accounts table: %s\n", err.Error())
		panic(msg)
	}
	_, err = s.DB.Exec(
		"create table if not exists orders(Number bigint NOT NULL UNIQUE, Status text NULL UNIQUE, Accural double precision, UploadedAt TIMESTAMP WITH TIME ZONE, Owner text)",
	)
	if err != nil {
		msg, _ := fmt.Printf("Failure during initialisation of orders table: %s\n", err.Error())
		panic(msg)
	}
	_, err = s.DB.Exec(
		"create table if not exists pending_orders(order_id text NOT NULL UNIQUE)",
	)
	if err != nil {
		msg, _ := fmt.Printf("Failure during initialisation of pending orders table: %s\n", err.Error())
		panic(msg)
	}
	_, err = s.DB.Exec(
		"create table if not exists withdrawns(OrderID bigint NOT NULL UNIQUE, Sum double precision, ProcessedAt TIMESTAMP WITH TIME ZONE, owner text)",
	)
	if err != nil {
		msg, _ := fmt.Printf("Failure during initialisation of withdrawns table: %s\n", err.Error())
		panic(msg)
	}
	logging.Debug(" SQL Repository initialisation successeeded")
	return s
}
