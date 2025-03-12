package main

import (
	
	"github.com/gin-gonic/gin"
	"github.com/Afomiat/ChatApp/infrastructure"
	"github.com/Afomiat/ChatApp/delivery/routers"



)

func main() {
	infrastructure.LoadEnv()
	mongoURI := infrastructure.GetEnv("MONGO_URI")

	db := infrastructure.ConnectMongo(mongoURI)

	r := gin.Default()
	routers.SetupRoutes(r, db)

	r.Run(":8080")
}