package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	custom_errors "github.com/slinkeres/segment-service/internal/errors"
	"github.com/slinkeres/segment-service/internal/models"
)

type UserSegmentRepository struct{
	db *pgxpool.Pool
}

func NewUserSegmentReposiotry(db *pgxpool.Pool) *UserSegmentRepository{
	return &UserSegmentRepository{
		db: db,
	}
}

func (r *UserSegmentRepository) AddUserToSegment(ctx context.Context, userSegment *models.UserSegment) error{
	query := `INSERT INTO user_segments VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userSegment.UserId, userSegment.SegmentId, userSegment.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr){
			if pgErr.Code == "23505"{
				return custom_errors.UserAlreadyExistsInSegment
			}
		}
		return err
	}
	return nil
}

func (r *UserSegmentRepository) DeleteUserFromSegment(ctx context.Context, userSegment *models.UserSegment) (int64, error){
	query := `DELETE FROM user_segments WHERE user_id = $1 AND segment_id = $2`
	rows, err := r.db.Exec(ctx, query, userSegment.UserId, userSegment.SegmentId)
	if err != nil {
		return 0, err
	}
	return rows.RowsAffected(), nil
}


func (r *UserSegmentRepository) GetSegmentsFromUser(ctx context.Context, id int) ([]models.Segment,error){
	query := `
		SELECT s.id, s.name, s.created_at
		FROM segments as s
		JOIN user_segments as us on s.id = us.segment_id
		JOIN users as u on u.id = us.user_id
		WHERE u.id = $1
	`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var segments []models.Segment
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