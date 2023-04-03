package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/madeinatria/love-all-backend/internal/controllers"
	"github.com/madeinatria/love-all-backend/internal/database"
)

func AllRoutes(route *gin.RouterGroup) {
	userController := controllers.NewUserController(database.Db)
	merchantController := controllers.NewMerchantController(database.Db)
	offerController := controllers.NewMerchantOfferController(database.Db)
	cardController := controllers.NewCardSubscriptionController(database.Db)
	transactionController := controllers.NewTransactionController(database.Db)
	// User resource
	route.GET("/users", userController.GetAllUsers)
	route.GET("/users/:id", userController.GetUser)
	route.POST("/users", userController.CreateUser)
	route.PUT("/users/:id", userController.UpdateUser)
	route.DELETE("/users/:id", userController.DeleteUser)

	// MerchantInfo
	route.GET("/merchants", merchantController.GetAllMerchants)
	route.GET("/merchants/:id", merchantController.GetMerchant)
	route.POST("/merchants", merchantController.CreateMerchant)
	route.PUT("/merchants/:id", merchantController.UpdateMerchant)
	route.DELETE("/merchants/:id", merchantController.DeleteMerchant)

	// Merchant Offer
	route.GET("/offers", offerController.GetAllMerchantOffers)
	route.GET("/offers/:id", offerController.GetMerchantOffer)
	route.POST("/offers", offerController.CreateMerchantOffer)
	route.PUT("/offers/:id", offerController.UpdateMerchantOffer)
	route.DELETE("/offers/:id", offerController.DeleteMerchantOffer)

	// Card Offer
	route.GET("/subscriptions", cardController.GetAllCardSubscriptions)
	route.GET("/subscriptions/:id", cardController.GetCardSubscription)
	route.POST("/subscriptions", cardController.CreateCardSubscription)
	route.PUT("/subscriptions/:id", cardController.UpdateCardSubscription)
	route.DELETE("/subscriptions/:id", cardController.DeleteCardSubscription)

	// Transaction
	route.GET("/transactions", transactionController.GetAllTransaction)
	route.GET("/transactions/:id", transactionController.GetTransaction)
	route.POST("/transactions", transactionController.CreateTransaction)
	route.PUT("/transactions/:id", transactionController.UpdateTransaction)
	route.DELETE("/transactions/:id", transactionController.DeleteTransaction)
}
