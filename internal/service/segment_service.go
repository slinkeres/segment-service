package service

import (
	"context"
	"time"

	"github.com/slinkeres/segment-service/internal/errors"
	"github.com/slinkeres/segment-service/internal/models"
	"github.com/slinkeres/segment-service/internal/repository/postgres"
)

type SegmentService struct{
	repo *postgres.SegmentRepository
}

func NewSegmentService(repo *postgres.SegmentRepository) *SegmentService{
	return &SegmentService{
		repo: repo,
	}
}

func (s *SegmentService) Create(ctx context.Context, segment *models.Segment) error{
	segment.CreatedAt = time.Now()
	err := s.repo.Create(ctx, segment)
	if err != nil {
		return err
	}
	return nil
}

func (s *SegmentService) GetAll(ctx context.Context) ([]models.Segment, error){
	segments, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return segments, nil
}

func (s *SegmentService) GetById(ctx context.Context, id int) (*models.Segment, error){
	segment, err := s.repo.GetById(ctx, id)
	if err != nil{
		return nil, err
	}
	return segment, nil
}

func (s *SegmentService) Update(ctx context.Context, id int, name string) error{
	if name == ""{
		return custom_errors.SegmentNameEmpty
	}
	rows, err := s.repo.Update(ctx, id, name)
	if err != nil {
		return err
	}
	if rows == 0{
		return custom_errors.SegmentNotFound
	}
	return nil
}

func (s *SegmentService) Delete(ctx context.Context, id int) (error){
	rows, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0{
		return custom_errors.SegmentNotFound
	}
	return nil
}