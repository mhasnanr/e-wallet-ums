package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/e-wallet/bootstrap"
)

func ServeHTTP() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "server is healthy"})
	})

	server := &http.Server{Addr: ":" + bootstrap.GetEnv("HTTP_PORT", "8080"), Handler: r}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server stopped")
	}
}
