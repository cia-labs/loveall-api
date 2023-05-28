package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeinatria/love-all-backend/internal/models"
	"github.com/madeinatria/love-all-backend/internal/utils"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{db}
}

func (tc *TransactionController) GetAllTransaction(c *gin.Context) {
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
	if err := tc.DB.Model(&models.Transaction{}).Count(&totalCount).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	var transactions []models.Transaction
	offset := (page - 1) * limit
	if err := tc.DB.Preload("CardSubscription.User").Preload("MerchantOffer.MerchantInfo.User").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	transactionResponses := make([]models.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = transaction.ToTransactionResponse()
	}

	c.JSON(200, gin.H{
		"data": transactionResponses,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalPages": utils.CalculateTotalPages(totalCount, int64(limit)),
			"totalCount": totalCount,
		},
	})
}

func (tc *TransactionController) GetTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var transaction models.Transaction
	err = tc.DB.Preload("CardSubscription.User").Preload("MerchantOffer.MerchantInfo.User").First(&transaction, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction.ToTransactionResponse())
}

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

	c.JSON(http.StatusOK, transaction.ToTransactionResponse())
}

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

	c.JSON(http.StatusOK, transaction.ToTransactionResponse())
}

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
