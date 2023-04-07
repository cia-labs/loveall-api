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

// GetAllMerchants returns a list of all merchants
// @Summary Get all merchants
// @Description Returns a list of all merchants in the system
// @Tags Merchants
// @Accept json
// @Produce json
// @Success 200 {array} models.MerchantInfo
// @Router /merchants [get]
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

// GetMerchant returns a single merchant by ID
// @Summary Get merchant by ID
// @Description Returns a single merchant by ID
// @Tags Merchants
// @Param id path int true "Merchant ID"
// @Accept json
// @Produce json
// @Success 200 {object} models.MerchantInfo
// @Router /merchants/{id} [get]
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

// GetMerchantsForUser returns all the merchant by userID
// @Summary Get merchants by userID
// @Description Returns merchants by userID
// @Tags Merchants
// @Param id path int true "User ID"
// @Accept json
// @Produce json
// @Success 200 {object} models.MerchantInfo
// @Router /merchantsbyuser/{id} [get]
func (mc *MerchantController) GetMerchantsForUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	var merchant []models.MerchantInfo
	err = mc.db.Preload("User").Where("user_id = ?", id).Find(&merchant).Error // find product with code D42
	// err = mc.db.Preload("User").First(&merchant, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": merchant,
	})
}

// CreateMerchant godoc
// @Summary Create a new merchant
// @Description Create a new merchant with the provided details
// @Tags merchants
// @Accept json
// @Produce json
// @Param Merchant body models.MerchantInfo true "Merchant details"
// @Success 201
// @Failure 400
// @Router /merchants [post]
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

// UpdateMerchant godoc
// @Summary Update an existing merchant
// @Description Update an existing merchant with the provided details
// @Tags merchants
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Param Merchant body models.MerchantInfo true "Merchant details"
// @Success 200
// @Failure 400
// @Router /merchants/{id} [put]
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
