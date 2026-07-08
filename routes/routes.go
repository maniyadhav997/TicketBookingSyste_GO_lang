package routes

import (
	"net/http"

	"ticket-system/handlers"
	"ticket-system/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the app.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "server is running",
		})
	})

	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)

	// Ticket routes (protected by JWT).
	r.POST("/tickets", middleware.AuthMiddleware(), handlers.CreateTicket)
	r.GET("/tickets", middleware.AuthMiddleware(), handlers.ListTickets)
	r.GET("/tickets/:id", middleware.AuthMiddleware(), handlers.GetTicketByID)
	r.PATCH("/tickets/:id/status", middleware.AuthMiddleware(), handlers.UpdateTicketStatus)
}

