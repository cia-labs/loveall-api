package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeinatria/love-all-backend/internal/models"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{db}
}

//	func (tc *TransactionController) GetAllTransaction(c *gin.Context) {
//		var tcSubs []models.Transaction
//		err := tc.DB.Find(&tcSubs).Error
//		if err != nil {
//			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(200, tcSubs)
//	}

// GetAllTransaction godoc
// @Summary Get all transactions
// @Description Get all transactions available
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Router /transactions [get]
func (tc *TransactionController) GetAllTransaction(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	var totalCount int64
	if err := tc.DB.Model(&models.Transaction{}).Count(&totalCount).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	var transactions []models.Transaction
	offset := (page - 1) * limit
	if err := tc.DB.Preload("User").Preload("CardSubscription.User").Preload("MerchantOffer.MerchantInfo.User").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": transactions,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalPages": int(math.Ceil(float64(totalCount) / float64(limit))),
			"totalCount": totalCount,
		},
	})
}

// GetTransaction godoc
// @Summary Get a specific transaction
// @Description Get a specific transaction by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200
// @Failure 400
// @Router /transactions/{id} [get]
func (tc *TransactionController) GetTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var transaction models.Transaction
	err = tc.DB.Preload("User").Preload("MerchantOffer").First(&transaction, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction with the provided details
// @Tags transactions
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Router /transactions [post]
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	err := c.ShouldBindJSON(&transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = tc.DB.Create(&transaction).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction godoc
// @Summary Update an existing transaction
// @Description Update a transaction with a specific ID
// @Tags transaction
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200
// @Failure 400
// @Router /transactions/{id} [put]
func (tc *TransactionController) UpdateTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var transaction models.Transaction
	err = tc.DB.First(&transaction, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var updatedTransaction models.Transaction
	err = c.ShouldBindJSON(&updatedTransaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = tc.DB.Model(&transaction).Updates(updatedTransaction).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction godoc
// @Summary Delete a transaction
// @Description Delete a transaction with a specific ID
// @Tags transaction
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 204 ""
// @Failure 400
// @Router /transactions/{id} [delete]
func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var transaction models.Transaction
	err = tc.DB.First(&transaction, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	err = tc.DB.Delete(&transaction).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
