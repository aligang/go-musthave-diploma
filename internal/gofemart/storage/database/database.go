package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"sync"
)
import _ "github.com/jackc/pgx/v4/stdlib"

type Storage struct {
	DB   *sql.DB
	Tx   map[context.Context]*sql.Tx
	Lock sync.Mutex
}

func New(conf *config.Config) *Storage {
	logging.Debug("Initialisating SQL Repository")
	db, err := sql.Open("pgx", conf.DatabaseURI)
	if err != nil {
		panic(err)
	}
	s := &Storage{
		DB: db,
		Tx: map[context.Context]*sql.Tx{},
	}
	_, err = s.DB.Exec(
		"create table if not exists accounts(Login text, Password text, Current double precision, Withdraw double precision)",
	)
	if err != nil {
		msg, _ := fmt.Printf("Failure during initialisation of accounts table: %s\n", err.Error())
		panic(msg)
	}
	_, err = s.DB.Exec(
		"create table if not exists orders(Number bigint, Status text, Accural double precision, UploadedAt TIMESTAMP WITH TIME ZONE, Owner text)",
	)
	if err != nil {
		msg, _ := fmt.Printf("Failure during initialisation of orders table: %s\n", err.Error())
		panic(msg)
	}
	_, err = s.DB.Exec(
		"create table if not exists pending_orders(order_id text)",
	)
	if err != nil {
		msg, _ := fmt.Printf("Failure during initialisation of pending orders table: %s\n", err.Error())
		panic(msg)
	}
	_, err = s.DB.Exec(
		"create table if not exists withdrawns(Order_id bigint, Sum double precision, Processed_At TIMESTAMP WITH TIME ZONE, owner text)",
	)
	if err != nil {
		msg, _ := fmt.Printf("Failure during initialisation of withdrawns table: %s\n", err.Error())
		panic(msg)
	}
	logging.Debug(" SQL Repository initialisation successeeded")
	return s
}
