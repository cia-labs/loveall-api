package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/madeinatria/love-all-backend/internal/models"
	"github.com/madeinatria/love-all-backend/internal/utils"
)

type MerchantOfferController struct {
	db *gorm.DB
}

func NewMerchantOfferController(db *gorm.DB) *MerchantOfferController {
	return &MerchantOfferController{db}
}

func (moc *MerchantOfferController) GetAllMerchantOffers(c *gin.Context) {
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
	if err := moc.db.Model(&models.MerchantOffer{}).Count(&totalCount).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var merchantOffers []models.MerchantOffer
	offset := (page - 1) * limit
	if err := moc.db.Preload("MerchantInfo").Offset(offset).Limit(limit).Find(&merchantOffers).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var merchantOfferResponses []models.MerchantOfferResponse
	for _, merchantOffer := range merchantOffers {
		merchantOfferResponses = append(merchantOfferResponses, merchantOffer.ToMerchantOfferResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": merchantOfferResponses,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalPages": utils.CalculateTotalPages(totalCount, int64(limit)),
			"totalCount": totalCount,
		},
	})
}

func (moc *MerchantOfferController) GetMerchantOffer(c *gin.Context) {
	id := c.Param("id")

	var merchantOffer models.MerchantOffer
	err := moc.db.Preload("MerchantInfo").First(&merchantOffer, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Merchant offer not found"})
		return
	}

	c.JSON(http.StatusOK, merchantOffer.ToMerchantOfferResponse())
}

func (moc *MerchantOfferController) CreateMerchantOffer(c *gin.Context) {
	var merchantOffer models.MerchantOffer
	if err := c.BindJSON(&merchantOffer); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := moc.db.Create(&merchantOffer).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, merchantOffer.ToMerchantOfferResponse())
}

func (moc *MerchantOfferController) UpdateMerchantOffer(c *gin.Context) {
	id := c.Param("id")

	var merchantOffer models.MerchantOffer
	err := moc.db.First(&merchantOffer, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Merchant offer not found"})
		return
	}

	if err := c.BindJSON(&merchantOffer); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := moc.db.Save(&merchantOffer).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, merchantOffer.ToMerchantOfferResponse())
}

func (moc *MerchantOfferController) DeleteMerchantOffer(c *gin.Context) {
	id := c.Param("id")

	var merchantOffer models.MerchantOffer
	err := moc.db.First(&merchantOffer, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Merchant offer not found"})
		return
	}

	if err := moc.db.Delete(&merchantOffer).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}
