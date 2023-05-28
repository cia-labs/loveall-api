package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeinatria/love-all-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db}
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	var totalCount int64
	if err := uc.db.Model(&models.User{}).Count(&totalCount).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userResponses []models.UserResponse
	offset := (page - 1) * limit
	if err := uc.db.Model(&models.User{}).Offset(offset).Limit(limit).Select("id, email, first_name, last_name").Find(&userResponses).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userResponses,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalPages": int(math.Ceil(float64(totalCount) / float64(limit))),
			"totalCount": totalCount,
		},
	})
}

func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.UserResponse
	err := uc.db.Model(&models.User{}).First(&user, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	// Hash the user's password using bcrypt.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error hashing password"})
		return
	}

	// Replace the user's plaintext password with the hashed version.
	user.Password = string(hashedPassword)

	if err := uc.db.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	err := uc.db.First(&user, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := uc.db.Save(&user).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}

func (uc *UserController) DeleteUser(c *gin.Context) {

	id := c.Param("id")
	var user models.User
	err := uc.db.First(&user, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := uc.db.Delete(&user).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}
