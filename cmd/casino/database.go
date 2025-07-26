package main

import "github.com/jackc/pgx/v5/pgxpool"

type Database struct {
	pgPool *pgxpool.Conn
}
