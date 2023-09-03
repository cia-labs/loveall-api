package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeinatria/love-all-backend/internal/models"
	"github.com/madeinatria/love-all-backend/internal/utils"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	var totalCount int64
	if err := uc.db.Model(&models.User{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var users []models.User
	offset := (page - 1) * limit
	if err := uc.db.Model(&models.User{}).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userResponses,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalPages": utils.CalculateTotalPages(totalCount, int64(limit)),
			"totalCount": totalCount,
		},
	})
}

func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := uc.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Println("DEBUG: Compare Hash error:", err.Error())
		return
	}
	c.JSON(http.StatusOK, user.ToResponse())
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Debug: Create user req body error=", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		log.Println("DEBUG: Error hashing password:", err.Error())
		return
	}

	// Convert the hashedPassword to string
	user.Password = string(hashedPassword)

	log.Println("INFO: Creating user!")
	if err := uc.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	log.Println("DEBUG: inside:")
	var user models.User
	if err := uc.db.First(&user, id).Error; err != nil {
		log.Println("DEBUG: Compare Hash error:", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("DEBUG: Compare Hash error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := uc.db.Save(&user).Error; err != nil {
		log.Println("DEBUG: Compare Hash error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := uc.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := uc.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
