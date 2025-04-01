package main

import (
	"context"
	"log"

	"github.com/Afomiat/ChatApp/delivery/routers"
	"github.com/Afomiat/ChatApp/infrastructure"
	"github.com/Afomiat/ChatApp/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// In main.go
func main() {
    infrastructure.LoadEnv()
    mongoURI := infrastructure.GetEnv("MONGO_URI")
    db := infrastructure.ConnectMongo(mongoURI)

    // Initialize OTP repository indexes
    otpRepo := repository.NewOtpRepository(db)
    if err := otpRepo.Initialize(context.Background()); err != nil {
        log.Fatal("Failed to create OTP indexes:", err)
    }

    r := gin.Default()
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:5173"},
        AllowMethods: []string{"GET", "POST", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
    }))

    routers.Setup(r, db)
    r.Run(":8080")
}