package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/slinkeres/segment-service/internal/service"
)

type UserHandler struct{
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler{
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Create(c *gin.Context){
	user,err := h.service.Create(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to create user in database",
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":"user created",
		"id":user.Id,
	})
}

func (h *UserHandler) GetAll(c *gin.Context){
	users,err := h.service.GetAll(c.Request.Context())
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed to get all users in database",
			"error":err.Error(),
		})
		return
	}
	if len(users) == 0{
		c.JSON(http.StatusOK, gin.H{
			"message":"users empty",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": users,
	})
}

func(h *UserHandler) GetById(c *gin.Context){
	id , err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get id",
			"error":err.Error(),
		})
		return
	}
	user, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get user by id ",
			"error":err.Error(),
		})
		return
	}
	if user == nil{
		c.JSON(http.StatusNotFound, gin.H{
			"message":"user not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result":user,
	})
}

func (h *UserHandler) Delete(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed get id",
			"error":err.Error(),
		})
		return
	}
	err = h.service.Delete(c.Request.Context(), id)
	if err != nil{
		if err.Error() == "user not found"{
			c.JSON(http.StatusNotFound, gin.H{
				"message":"user not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message":"failed delete user by id",
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":"user deleted",
		"id":id,
	})
}