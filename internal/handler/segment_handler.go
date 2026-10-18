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

type SegmentHandler struct {
	service *service.SegmentService
}

func NewSegmentHandler(service *service.SegmentService) *SegmentHandler{
	return &SegmentHandler{
		service: service,
	}
}

func (h *SegmentHandler) Create(c *gin.Context){
	var segment models.Segment
	err := c.ShouldBindJSON(&segment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get segment parameters",
			"error":err.Error(),
		})
		return
	}
	err = h.service.Create(c.Request.Context(), &segment)
	if err != nil {
		if errors.Is(err, custom_errors.SegmentNameExists){
			c.JSON(http.StatusConflict, gin.H{
				"message":"segment name exists",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed create segment",
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":"segment created",
		"id": segment.Id,
		"name": segment.Name,
	})
}

func (h *SegmentHandler) GetAll(c *gin.Context){
	segments, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get all segments",
			"error":err.Error(),
		})
		return
	}
	if len(segments) == 0{
		c.JSON(http.StatusOK, gin.H{
			"message":"segments empty",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": segments,
	})
}

func (h *SegmentHandler) GetById(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get id",
			"error":err.Error(),
		})
		return
	}
	segment, err := h.service.GetById(c.Request.Context(), id)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get segment by id",
			"error":err.Error(),
		})
		return
	}
	if segment == nil{
		c.JSON(http.StatusNotFound, gin.H{
			"message":"segment not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result":segment,
	})
}

func (h *SegmentHandler) Update(c *gin.Context){
	id, err:= strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get id",
			"error":err.Error(),
		})
		return
	}
	var segment models.Segment
	err = c.ShouldBindJSON(&segment)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get new segment parameters",
			"error":err.Error(),
		})
		return
	}
	err = h.service.Update(c.Request.Context(), id, segment.Name)
	if err != nil{
		if errors.Is(err, custom_errors.SegmentNameEmpty){
			c.JSON(http.StatusBadRequest, gin.H{
				"message":"segment cannot be empty",
			})
			return
		} 
		if errors.Is(err, custom_errors.SegmentNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"segment not found",
			})
			return
		}
		if errors.Is(err, custom_errors.SegmentNameExists){
			c.JSON(http.StatusConflict, gin.H{
				"message":"segment name exists",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed update segment",
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":"segment updated",
	})
}

func (h *SegmentHandler) Delete(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get id",
			"error":err.Error(),
		})
		return
	}
	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, custom_errors.SegmentNotFound){
			c.JSON(http.StatusNotFound, gin.H{
				"message":"segment not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed delete segment",
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":"segment deleted",
		"id":id,
	})
}