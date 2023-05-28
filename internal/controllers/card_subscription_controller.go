package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/madeinatria/love-all-backend/internal/database"
	"github.com/madeinatria/love-all-backend/internal/models"
	"github.com/madeinatria/love-all-backend/internal/utils"
)

type CardSubscriptionController struct {
	db *gorm.DB
}

func NewCardSubscriptionController(db *gorm.DB) *CardSubscriptionController {
	return &CardSubscriptionController{
		db: db}
}

func (csc *CardSubscriptionController) GetAllCardSubscriptions(c *gin.Context) {
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
	if err := csc.db.Preload("User").Model(&models.CardSubscription{}).Count(&totalCount).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var cardSubs []models.CardSubscription
	offset := (page - 1) * limit
	if err := csc.db.Preload("User").Offset(offset).Limit(limit).Find(&cardSubs).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cardSubResponses := make([]models.CardSubscriptionResponse, len(cardSubs))
	for i, cardSub := range cardSubs {
		cardSubResponses[i] = cardSub.ToCardSubscriptionResponse()
	}

	c.JSON(200, gin.H{
		"data": cardSubResponses,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalPages": utils.CalculateTotalPages(totalCount, int64(limit)),
			"totalCount": totalCount,
		},
	})
}

func (csc *CardSubscriptionController) GetCardSubscription(c *gin.Context) {
	id := c.Param("id")
	var cardSub models.CardSubscription
	err := csc.db.Preload("User").First(&cardSub, id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Card subscription not found"})
		return
	}
	c.JSON(http.StatusOK, cardSub.ToCardSubscriptionResponse())
}

func (csc *CardSubscriptionController) CreateCardSubscription(c *gin.Context) {
	var cardSub models.CardSubscription
	if err := c.BindJSON(&cardSub); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := csc.db.Create(&cardSub).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, cardSub.ToCardSubscriptionResponse())
}

func (csc *CardSubscriptionController) ValidateCardSubscription(c *gin.Context) {
	var cardValidate models.ValidateRequest
	if err := c.BindJSON(&cardValidate); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var cardSubscription models.CardSubscription
	cardErr := database.Db.Preload("User").Where("id = ?", cardValidate.CardId).First(&cardSubscription)
	if cardErr.Error != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Card subscription not found"})
		return
	}
	// Check if the merchant has any avaiable or valid card discount for this user.
	var merchantInfo models.MerchantInfo
	merchErr := database.Db.Where("id = ?", cardValidate.MerchantId).First(&merchantInfo)

	if merchErr.Error != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "merchant not found"})
		return
	}

	var matchingOffer models.MerchantOffer
	offerErr := database.Db.Preload("MerchantInfo.User").Where("merchant_info_id = ? AND card_name = ?",
		cardValidate.MerchantId,
		cardSubscription.CardName).Find(&matchingOffer)
	if offerErr.Error != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "no valid offer found for the merchant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"card":  cardSubscription,
		"offer": matchingOffer,
	})

}

func (csc *CardSubscriptionController) UpdateCardSubscription(c *gin.Context) {
	id := c.Param("id")
	var cardSub models.CardSubscription
	err := csc.db.First(&cardSub, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Card subscription not found"})
		return
	}

	if err := c.BindJSON(&cardSub); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := csc.db.Save(&cardSub).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, cardSub.ToCardSubscriptionResponse())
}

func (csc *CardSubscriptionController) DeleteCardSubscription(c *gin.Context) {
	id := c.Param("id")

	var cardSub models.CardSubscription
	err := csc.db.First(&cardSub, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Card subscription not found"})
		return
	}

	if err := csc.db.Delete(&cardSub).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}
