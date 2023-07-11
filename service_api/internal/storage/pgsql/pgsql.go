package pgsql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"service_api/internal/config"
	"service_api/internal/metrics"
	"service_api/internal/storage"
	"time"
)

type PostgresStorage struct {
	conn *pgxpool.Pool
}

// NewPostgresStorage метод создания нового экземпляра storage
func NewPostgresStorage(ctx context.Context, cfg config.DBSQLConnection) (*PostgresStorage, error) {
	connectConfig, err := pgxpool.ParseConfig(cfg.Uri)
	if err != nil {
		return nil, err
	}
	connectConfig.MaxConns = int32(cfg.MaxOpenConn)
	connectConfig.MaxConnIdleTime = time.Minute * time.Duration(cfg.ConnMaxLifetimeMinute)
	pool, err := pgxpool.NewWithConfig(ctx, connectConfig)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{
		conn: pool,
	}, nil
}

func (p *PostgresStorage) GetProduct(ctx context.Context, limit int, offset int) ([]storage.Product, error) {
	metricsObservePostgres := metrics.ObservePostgresQueryDuration()
	rows, err := p.conn.Query(ctx, "SELECT id, price FROM product LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []storage.Product
	for rows.Next() {
		var data storage.Product
		err := rows.Scan(&data.Id, &data.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, data)
	}
	metricsObservePostgres()
	return products, nil
}
