package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/slinkeres/segment-service/internal/models"
)

type UserRepository struct{
	db *pgxpool.Pool
}

func NewUserReposiotry(db *pgxpool.Pool) *UserRepository{
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error{
	query := `INSERT INTO users (created_at) VALUES($1) RETURNING id`
	err := r.db.QueryRow(ctx, query, user.CreatedAt).Scan(&user.Id)
	if err != nil {
		return err
	}
	return nil
}


func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error){
	query := `SELECT * FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next(){
		var user models.User
		err = rows.Scan(&user.Id, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetById(ctx context.Context, id int) (*models.User, error){
	query := "SELECT * FROM users WHERE id = $1"
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.Id, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows){
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Delete(ctx context.Context, id int) (int64, error){
	query := "DELETE FROM users WHERE id = $1"
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return 0,err
	}
	return res.RowsAffected(), nil
}