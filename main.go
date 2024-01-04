package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Use(requestLogger)
	InitializeRoutes(e)
	e.Start(":5050")
}

func requestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Record the start time of the request
		start := time.Now()

		// Call the next middleware or handler
		err := next(c)

		// Record the end time of the request
		end := time.Now()

		// Calculate the response time
		responseTime := end.Sub(start)

		// Log the request information
		log.Printf("[%s] %s - %s - %v", c.Request().Method, c.Path(), c.Response().Status, responseTime)

		return err
	}
}
