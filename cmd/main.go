package main

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/slinkeres/segment-service/config"
	"github.com/slinkeres/segment-service/internal/handler"
	"github.com/slinkeres/segment-service/internal/repository/postgres"
	"github.com/slinkeres/segment-service/internal/service"
)


func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil{
		slog.Error("failed load .env", "error", err)
		return
	}
	db, err := postgres.CreateDb(ctx, cfg.DBURL)
	if err != nil{
		slog.Error("Cannot create db", "error", err)
		return 
	}
	defer db.Close()
	slog.Info("db created")
	userRepository := postgres.NewUserReposiotry(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	segmentRepository := postgres.NewSegmentRepository(db)
	segmentService := service.NewSegmentService(segmentRepository)
	segmentHandler := handler.NewSegmentHandler(segmentService)

	userSegmentReposiotry := postgres.NewUserSegmentReposiotry(db)
	userSegmentService := service.NewUserSegmentService(userRepository, segmentRepository, userSegmentReposiotry)
	userSegmentHandler := handler.NewUserSegmentHandler(userSegmentService)

	router := gin.Default()

	router.POST("/users", userHandler.Create)
	router.GET("/users", userHandler.GetAll)
	router.GET("/users/:id", userHandler.GetById)
	router.DELETE("/users/:id", userHandler.Delete)

	router.POST("/segments", segmentHandler.Create)
	router.GET("/segments", segmentHandler.GetAll)
	router.GET("/segments/:id", segmentHandler.GetById)
	router.PATCH("/segments/:id", segmentHandler.Update)
	router.DELETE("/segments/:id", segmentHandler.Delete)

	router.POST("/user-segments/:segment_id/users/:user_id", userSegmentHandler.AddUserToSegment)
	router.DELETE("/user-segments/:segment_id/users/:user_id", userSegmentHandler.DeleteUserFromSegment)
	router.GET("/users/:id/segments", userSegmentHandler.GetSegmentsFromUser)

	router.POST("/segments/:id/apply-percent", userSegmentHandler.ApplyPercentToSegment)

	slog.Info("server starting at localhost:8080")
	router.Run(":"+cfg.Port)
}