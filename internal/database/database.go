package database

import (
	"fmt"
	"log"

	"github.com/madeinatria/love-all-backend/internal/models"
	// "github.com/mgechev/revive/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB
var DbErr error

func init() {

	dbUri := fmt.Sprintf("nimbus.db")
	Db, DbErr = gorm.Open(sqlite.Open(dbUri), &gorm.Config{})
	if DbErr != nil {
		log.Fatal("Could not connect to the database:", DbErr.Error())
	}

	// Automigrate database tables
	Db.AutoMigrate(&models.User{}, &models.CardSubscription{}, &models.MerchantOffer{}, &models.Transaction{}, &models.MerchantInfo{})
}
