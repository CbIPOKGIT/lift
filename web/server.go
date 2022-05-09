package web

import "github.com/gin-gonic/gin"

func StartServer() *gin.Engine {
	server := gin.Default()
	return server
}
