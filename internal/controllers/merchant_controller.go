package controllers

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/madeinatria/love-all-backend/internal/models"
)

type MerchantController struct {
	db *gorm.DB
}

func NewMerchantController(db *gorm.DB) *MerchantController {
	return &MerchantController{db}
}

// func (mc *MerchantController) GetAllMerchants(c *gin.Context) {
// 	var merchants []models.MerchantInfo
// 	err := mc.db.Preload("User").Find(&merchants).Error
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, merchants)
// }

func (mc *MerchantController) GetAllMerchants(c *gin.Context) {
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
	if err := mc.db.Model(&models.MerchantInfo{}).Count(&totalCount).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var merchants []models.MerchantInfo
	offset := (page - 1) * limit
	if err := mc.db.Preload("User").Offset(offset).Limit(limit).Find(&merchants).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": merchants,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalPages": int(math.Ceil(float64(totalCount) / float64(limit))),
			"totalCount": totalCount,
		},
	})
}

func (mc *MerchantController) GetMerchant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID"})
		return
	}

	var merchant models.MerchantInfo
	err = mc.db.Preload("User").First(&merchant, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, merchant)
}

func (mc *MerchantController) CreateMerchant(c *gin.Context) {
	var merchant models.MerchantInfo
	if err := c.BindJSON(&merchant); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := mc.db.Create(&merchant).Error
	if err != nil {
		// if isDuplicateKeyError(err) {
		// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Merchant with this name already exists"})
		// } else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// }
		return
	}

	c.JSON(http.StatusCreated, merchant)
}

func (mc *MerchantController) UpdateMerchant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID"})
		return
	}

	var merchant models.MerchantInfo
	err = mc.db.First(&merchant, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := c.BindJSON(&merchant); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = mc.db.Save(&merchant).Error
	if err != nil {
		// if isDuplicateKeyError(err) {
		// c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Merchant with this name already exists"})
		// } else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// }
		return
	}

	c.JSON(http.StatusOK, merchant)
}

func (mc *MerchantController) DeleteMerchant(c *gin.Context) {
	//TODO: have to implement delete user
}
