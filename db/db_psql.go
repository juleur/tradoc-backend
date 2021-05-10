package db

import "github.com/jackc/pgx/v4/pgxpool"

type DBPsql struct {
	db *pgxpool.Pool
}

func NewDBPsql(db *pgxpool.Pool) *DBPsql {
	return &DBPsql{db: db}
}
