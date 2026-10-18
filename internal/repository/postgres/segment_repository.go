package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	custom_errors "github.com/slinkeres/segment-service/internal/errors"
	"github.com/slinkeres/segment-service/internal/models"
)

type SegmentRepository struct {
	db *pgxpool.Pool
}

func NewSegmentRepository(db *pgxpool.Pool) *SegmentRepository {
	return &SegmentRepository{
		db: db,
	}
}

func (r *SegmentRepository) Create(ctx context.Context, segment *models.Segment) error{
	query := `INSERT INTO segments (name, created_at) VALUES($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, segment.Name, segment.CreatedAt).Scan(&segment.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr){
			if pgErr.Code == "23505"{
				return custom_errors.SegmentNameExists
			}
		}
		return err
	}
	return nil
}

func (r *SegmentRepository) GetAll(ctx context.Context) ([]models.Segment, error){
	query := `SELECT * FROM segments`
	var segments []models.Segment
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var segment models.Segment
		err = rows.Scan(&segment.Id, &segment.Name, &segment.CreatedAt)
		if err != nil{
			return nil, err
		}
		segments = append(segments, segment)
	}
	err = rows.Err()
	if err != nil{
		return nil, err
	}
	return segments, nil
}

func (r *SegmentRepository) GetById(ctx context.Context, id int) (*models.Segment, error){
	query := `SELECT * FROM segments WHERE id = $1`
	var segment models.Segment
	err := r.db.QueryRow(ctx, query, id).Scan(&segment.Id, &segment.Name, &segment.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows){
		return nil, nil
	}
	return &segment, err
}

func (r *SegmentRepository) Update(ctx context.Context, id int, name string) (int64, error){
	query := `UPDATE segments  SET  name = $1 WHERE id = $2`
	res, err := r.db.Exec(ctx, query, name, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr){
			if pgErr.Code == "23505"{
				return 0, custom_errors.SegmentNameExists
			}
		}
		return 0, err
	}
	return res.RowsAffected(), nil 
}

func (r *SegmentRepository) Delete(ctx context.Context, id int) (int64, error){
	query := `DELETE FROM segments WHERE id = $1`
	rows, err := r.db.Exec(ctx, query, id)
	if err != nil{
		return 0, err
	}
	return rows.RowsAffected(), nil
}