package server

import (
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server/handlers"
	"github.com/gin-gonic/gin"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", handlers.Health)
	router.GET("/supported-languages", handlers.SupportedLanguages)
	router.PUT("/set-language", handlers.SetLanguagePreference)

	handlers.MetricsHandler(router)

	return router
}
