package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/madeinatria/love-all-backend/internal/models"
)

type CardSubscriptionController struct {
	db *gorm.DB
}

func NewCardSubscriptionController(db *gorm.DB) *CardSubscriptionController {
	return &CardSubscriptionController{db}
}

//	func (csc *CardSubscriptionController) GetAllCardSubscriptions(c *gin.Context) {
//		var cardSubs []models.CardSubscription
//		err := csc.db.Find(&cardSubs).Error
//		if err != nil {
//			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(200, cardSubs)
//	}

// GetAllCardSubscriptions godoc
// @Summary Get all card subscriptions
// @Description Get all card subscriptions available
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Router /subscriptions [get]
func (csc *CardSubscriptionController) GetAllCardSubscriptions(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	var totalCount int64
	if err := csc.db.Preload("User").Model(&models.CardSubscription{}).Count(&totalCount).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	var cardSubs []models.CardSubscription
	offset := (page - 1) * limit
	if err := csc.db.Preload("User").Offset(offset).Limit(limit).Find(&cardSubs).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": cardSubs,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalPages": int(math.Ceil(float64(totalCount) / float64(limit))),
			"totalCount": totalCount,
		},
	})
}

// GetCardSubscription godoc
// @Summary Get a specific card subscription
// @Description Get a specific card subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200
// @Failure 400
// @Router /subscriptions/{id} [get]
func (csc *CardSubscriptionController) GetCardSubscription(c *gin.Context) {
	id := c.Param("id")
	var cardSub models.CardSubscription
	err := csc.db.Preload("User").First(&cardSub, id).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Card subscription not found"})
		return
	}
	c.JSON(200, cardSub)
}

// CreateCardSubscription godoc
// @Summary Create a new card subscription
// @Description Create a new card subscription with the provided details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Router /subscriptions [post]
func (csc *CardSubscriptionController) CreateCardSubscription(c *gin.Context) {
	var cardSub models.CardSubscription
	if err := c.BindJSON(&cardSub); err != nil {
		// c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := csc.db.Create(&cardSub).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, cardSub)
}

// UpdateCardSubscription godoc
// @Summary Update an existing card subscription
// @Description Update an existing card subscription with the provided details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200
// @Failure 400
// @Router /subscriptions/{id} [put]
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

	c.JSON(200, cardSub)
}

// DeleteCardSubscription godoc
// @Summary Delete an existing card subscription
// @Description Delete an existing card subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 400
// @Router /subscriptions/{id} [delete]
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
