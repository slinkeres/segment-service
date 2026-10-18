package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	custom_errors "github.com/slinkeres/segment-service/internal/errors"
	"github.com/slinkeres/segment-service/internal/models"
	"github.com/slinkeres/segment-service/internal/service"
)

type UserSegmentHandler struct {
	service *service.UserSegmentService
}

func NewUserSegmentHandler(service *service.UserSegmentService) *UserSegmentHandler{
	return &UserSegmentHandler{
		service: service,
	}
}

func (h *UserSegmentHandler) AddUserToSegment(c *gin.Context){
	segmentId, err := strconv.Atoi(c.Param("segment_id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get segment id",
			"error":err.Error(),
		})
		return
	}
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get user id",
			"error":err.Error(),
		})
		return
	}
	var userSegment models.UserSegment
	userSegment.UserId = userId
	userSegment.SegmentId = segmentId
	err = h.service.AddUserToSegment(c.Request.Context(), &userSegment)
	if err != nil {
		if errors.Is(err, custom_errors.UserNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"user not found",
			})
			return
		}
		if errors.Is(err, custom_errors.SegmentNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"segment not found",
			})
			return
		}
		if errors.Is(err, custom_errors.UserAlreadyExistsInSegment){
			c.JSON(http.StatusConflict, gin.H{
				"message":"user already exists in segment",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed add user to segment",
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":"user add to segment",
	})
}

func (h *UserSegmentHandler) DeleteUserFromSegment(c *gin.Context){
	segmentId, err := strconv.Atoi(c.Param("segment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get segment id",
			"error":err.Error(),
		})
		return
	}
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get user id from segment",
			"error":err.Error(),
		})
		return
	}
	var userSegment models.UserSegment
	userSegment.UserId = userId
	userSegment.SegmentId = segmentId
	err = h.service.DeleteUserFromSegment(c.Request.Context(), &userSegment)
	if err != nil{
		if errors.Is(err, custom_errors.UserNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"user not found",
			})
			return
		}
		if errors.Is(err, custom_errors.SegmentNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"segment not found",
			})
			return
		}
		if errors.Is(err, custom_errors.UserSegmentNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"user not found in segment",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed delete user from segment",
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":"user deleted from segment",
	})
}

func (h *UserSegmentHandler) GetSegmentsFromUser(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get user id",
			"error":err.Error(),
		})
		return
	}
	segments, err := h.service.GetSegmentsFromUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, custom_errors.UserNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"user not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get segments from user",
			"error":err.Error(),
		})
		return
	}
	if segments == nil{
		c.JSON(http.StatusOK, gin.H{
			"result":"segments empty",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result":segments,
	})
}

func (h *UserSegmentHandler) ApplyPercentToSegment(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get id",
			"error":err.Error(),
		})
		return
	}
	var percent models.Percent
	err = c.ShouldBindJSON(&percent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get percent parameter",
			"error":err.Error(),
		})
		return
	}
	if percent.Num > 100 || percent.Num < 1{
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"percent musbe between 1 and 100",
		})
		return
	}
	res, err := h.service.ApplyPercentToSegment(c.Request.Context(), id, &percent)
	if err != nil {
		if errors.Is(err, custom_errors.SegmentNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"segment not found",
			})
			return
		}
		if errors.Is(err, custom_errors.UsersArrayEmpty){
			c.JSON(http.StatusNoContent, gin.H{
				"message":"users array empty",
			})
			return
		}
		if errors.Is(err, custom_errors.ZeroUsers){
			c.JSON(http.StatusNoContent, gin.H{
				"message":"zero users to add",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed add users to segment",
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":"users was added to segment",
		"result": res,
	})
}