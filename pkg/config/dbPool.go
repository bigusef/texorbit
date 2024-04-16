package config

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"time"
)

func initConfiguration(dbURL string) (*pgxpool.Config, error) {
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = int32(4)
	dbConfig.MinConns = int32(0)
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = time.Minute * 30
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		slog.Info("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		slog.Info("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		slog.Info("Closed the connection pool to the database!!")
	}

	return dbConfig, nil
}

func NewConnectionPool(ctx context.Context, dbURL string) *pgxpool.Pool {
	dbConfig, err := initConfiguration(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	connPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	connection, err := connPool.Acquire(ctx)
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	}
	defer connection.Release()

	err = connection.Ping(ctx)
	if err != nil {
		log.Fatal("Could not ping database")
	}

	slog.Info("Connected to the database!!")
	return connPool
}
