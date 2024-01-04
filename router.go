package main

import (
	"aviator/controller"

	"github.com/labstack/echo/v4"
)

// InitializeRoutes sets up all the routes for the application
func InitializeRoutes(e *echo.Echo) {
	apiGroup := e.Group("/api/v1")
	apiGroup.POST("/place_bet", controller.PlaceBet)
	apiGroup.POST("/cashout", controller.CashOut)
}
