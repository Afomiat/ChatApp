package main

import (
    "github.com/gin-gonic/gin"
    "github.com/Afomiat/ChatApp/infrastructure"
    "github.com/Afomiat/ChatApp/delivery/routers"
    "github.com/gin-contrib/cors"
)

func main() {
    infrastructure.LoadEnv()
    mongoURI := infrastructure.GetEnv("MONGO_URI")

    db := infrastructure.ConnectMongo(mongoURI)

    r := gin.Default()

    // Enable CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, // Allow frontend origin
        AllowMethods:     []string{"GET", "POST", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        // AllowWebSockets:  true,

    }))

    // Setup routes
    routers.SetupChatRoutes(r, db)
    routers.SetupUserRoutes(r, db)
    // connManager.StartCleanupRoutine(1 * time.Minute) // Checks every minute

    r.Run(":8080")
}