package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/external"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/handler"
	"github.com/mhasnanr/ewallet-ums/internal/middleware"
	"github.com/mhasnanr/ewallet-ums/internal/repository"
	"github.com/mhasnanr/ewallet-ums/internal/services"
	"gorm.io/gorm"
)

func ServeHTTP(db *gorm.DB) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "server is healthy"})
	})

	jwtManager := &helpers.JWTManager{}
	passwordHasher := &helpers.PasswordHasher{}
	walletExternal := &external.ExternalWallet{}

	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)

	authMiddleware := middleware.NewAuthMiddleware(sessionRepository, jwtManager)
	userService := services.NewUserService(userRepository, sessionRepository, jwtManager, passwordHasher, walletExternal)
	userHandler := handler.NewUserHandler(userService, authMiddleware)

	userHandler.RegisterRoute(r)

	server := &http.Server{Addr: ":" + bootstrap.GetEnv("HTTP_PORT", "8080"), Handler: r}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server stopped")
	}
}
