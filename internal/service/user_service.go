package service

import (
	"context"
	"time"

	custom_errors "github.com/slinkeres/segment-service/internal/errors"
	"github.com/slinkeres/segment-service/internal/models"
	"github.com/slinkeres/segment-service/internal/repository/postgres"
)


type UserService struct{
	repo *postgres.UserRepository
}

func NewUserService(repo *postgres.UserRepository) *UserService{
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context) (*models.User, error){
	user := &models.User{
		CreatedAt: time.Now(),
	}
	err := s.repo.Create(ctx,user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error){
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil 
}

func (s *UserService) GetById(ctx context.Context, id int) (*models.User, error){
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id int) error{
	rows, err := s.repo.Delete(ctx, id)
	if err != nil{
		return err
	}
	if rows == 0{
		return custom_errors.UserNotFound
	}
	return nil
}