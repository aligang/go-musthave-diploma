package database

import (
	"database/sql"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/logging"
)
import _ "github.com/jackc/pgx/v4/stdlib"

type DBStorage struct {
	DB *sql.DB
}

func New(conf *config.Config) *DBStorage {
	logging.Debug("Initialisating SQL Repository")
	db, err := sql.Open("pgx", conf.DatabaseURI)
	if err != nil {
		panic(err)
	}
	s := &DBStorage{
		DB: db,
	}
	rows, err := s.DB.Query(
		"create table if not exists metrics(ID text , MType text, Delta bigint, Value double precision, Hash text)",
	)
	if err != nil {

		panic(err.Error())
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	logging.Debug(" SQL Repository initialisation succesadead")
	return s
}
