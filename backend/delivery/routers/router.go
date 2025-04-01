package routers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(r *gin.Engine, db *mongo.Database) {
	NewChatRouter(r, db)
    NewUserRouter(r, db)
	NewRegisterRouter(r, db)
}