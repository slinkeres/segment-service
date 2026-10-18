package service

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"

	custom_errors "github.com/slinkeres/segment-service/internal/errors"
	"github.com/slinkeres/segment-service/internal/models"
	"github.com/slinkeres/segment-service/internal/repository/postgres"
)

type UserSegmentService struct{
	userRepo		*postgres.UserRepository
	segmentRepo		*postgres.SegmentRepository
	userSegmentRepo *postgres.UserSegmentRepository
}

func NewUserSegmentService(userRepo *postgres.UserRepository, segmentRepo *postgres.SegmentRepository, userSegmentRepo *postgres.UserSegmentRepository) *UserSegmentService{
	return &UserSegmentService{
		userRepo: userRepo,
		segmentRepo: segmentRepo,
		userSegmentRepo: userSegmentRepo,
	}
}

func (s *UserSegmentService) AddUserToSegment(ctx context.Context, userSegment *models.UserSegment) error{
	user, err := s.userRepo.GetById(ctx, userSegment.UserId)
	if err != nil{
		return err
	}
	if user == nil {
		return custom_errors.UserNotFound
	}
	segment, err := s.segmentRepo.GetById(ctx, userSegment.SegmentId)
	if err != nil{
		return err
	}
	if segment == nil {
		return custom_errors.SegmentNotFound
	}
	userSegment.CreatedAt = time.Now()
	err = s.userSegmentRepo.AddUserToSegment(ctx, userSegment)
	if err != nil{
		return err
	}
	return nil 
}

func (s *UserSegmentService) DeleteUserFromSegment(ctx context.Context, userSegment *models.UserSegment) error{
	user, err := s.userRepo.GetById(ctx, userSegment.UserId)
	if err != nil{
		return err
	}
	if user == nil{
		return custom_errors.UserNotFound
	}
	segment, err := s.segmentRepo.GetById(ctx, userSegment.SegmentId)
	if err != nil{
		return err
	}
	if segment == nil{
		return custom_errors.SegmentNotFound
	}
	rows, err := s.userSegmentRepo.DeleteUserFromSegment(ctx, userSegment)
	if err != nil {
		return err
	}
	if rows == 0{
		return custom_errors.UserSegmentNotFound
	}
	return nil
}

func (s *UserSegmentService) GetSegmentsFromUser(ctx context.Context, id int) ([]models.Segment, error){
	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil{
		return nil, custom_errors.UserNotFound
	}
	segments, err := s.userSegmentRepo.GetSegmentsFromUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return segments, nil
}

func (s *UserSegmentService) ApplyPercentToSegment(ctx context.Context, id int, percent *models.Percent) ([]models.User, error){
	segment, err := s.segmentRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if segment == nil{
		return nil, custom_errors.SegmentNotFound
	}
	users, err := s.userRepo.GetAll(ctx)
	if err != nil{
		return nil, err
	}
	if len(users) == 0{
		return nil, custom_errors.UsersArrayEmpty
	}
	count := int(math.Round(float64(len(users)) * float64(percent.Num)/100.0))
	if count == 0{
		return nil, custom_errors.ZeroUsers
	}
	rand.Shuffle(len(users), func(i, j int){
		users[i], users[j] = users[j], users[i]
	})
	added := 0
	var addedUser []models.User
	for _, user := range users{
		if added >= count {
			break
		}
		userSegment := models.UserSegment{
			UserId: user.Id,
			SegmentId: id,
			CreatedAt: time.Now(),
		}
		err = s.userSegmentRepo.AddUserToSegment(ctx, &userSegment)
		if err != nil {
			if errors.Is(err, custom_errors.UserAlreadyExistsInSegment){
				continue
			}
			return nil, err
		}
		addedUser = append(addedUser, user) 
		added++
	}
	return addedUser, nil
}