package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)



func CreateDb(ctx context.Context, dburl string) (*pgxpool.Pool, error){
	pool, err := pgxpool.New(ctx, dburl)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil{
		return nil, err
	}
	return pool, nil
}