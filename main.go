package main

import (
	"log"

	"github.com/Afomiat/Chatbot/Chatbot/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	server.SetupRoutes(r) 

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
