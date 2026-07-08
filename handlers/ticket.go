package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"ticket-system/database"
	"ticket-system/models"

	"github.com/gin-gonic/gin"
)

type CreateTicketInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}


func CreateTicket(c *gin.Context) {
	
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authorized",
		})
		return
	}

	var userID int
	switch v := userIDValue.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user id in token",
		})
		return
	}

	var input CreateTicketInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "title is required",
		})
		return
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)
	status := "open"

	query := `
		INSERT INTO tickets (title, description, status, user_id, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := database.DB.Exec(query, input.Title, input.Description, status, userID, createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not create ticket",
		})
		return
	}

	ticketID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ticket created but id not found",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ticket created successfully",
		"ticket": gin.H{
			"id":          ticketID,
			"title":       input.Title,
			"description": input.Description,
			"status":      status,
			"user_id":     userID,
			"created_at":  createdAt,
		},
	})
}


func ListTickets(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authorized",
		})
		return
	}

	var userID int
	switch v := userIDValue.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user id in token",
		})
		return
	}

	query := `
		SELECT id, title, description, status, user_id, created_at
		FROM tickets
		WHERE user_id = ?
		ORDER BY id DESC
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not fetch tickets",
		})
		return
	}
	defer rows.Close()

	tickets := []models.Ticket{}

	for rows.Next() {
		var t models.Ticket
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.UserID, &t.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not read ticket row",
			})
			return
		}

		tickets = append(tickets, t)
	}

	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not read tickets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "tickets fetched successfully",
		"tickets":  tickets,
	})
}

func GetTicketByID(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authorized",
		})
		return
	}

	var userID int
	switch v := userIDValue.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user id in token",
		})
		return
	}

	idParam := c.Param("id")
	ticketID, err := strconv.Atoi(idParam)
	if err != nil || ticketID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ticket id",
		})
		return
	}

	query := `
		SELECT id, title, description, status, user_id, created_at
		FROM tickets
		WHERE id = ? AND user_id = ?
	`

	var t models.Ticket
	err = database.DB.QueryRow(query, ticketID, userID).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.UserID,
		&t.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "ticket not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not fetch ticket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ticket fetched successfully",
		"ticket":  t,
	})
}

type UpdateTicketStatusInput struct {
	Status string `json:"status"`
}


func UpdateTicketStatus(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authorized",
		})
		return
	}

	var userID int
	switch v := userIDValue.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user id in token",
		})
		return
	}

	idParam := c.Param("id")
	ticketID, err := strconv.Atoi(idParam)
	if err != nil || ticketID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ticket id",
		})
		return
	}

	var input UpdateTicketStatusInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	if input.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "status is required",
		})
		return
	}

	allowedNext := map[string]map[string]bool{
		"open": {
			"in_progress": true,
			"open":        true,
		},
		"in_progress": {
			"closed":      true,
			"in_progress": true,
		},
		"closed": {
			"closed": true,
		},
	}

	var t models.Ticket
	selectQuery := `
		SELECT id, title, description, status, user_id, created_at
		FROM tickets
		WHERE id = ? AND user_id = ?
	`

	err = database.DB.QueryRow(selectQuery, ticketID, userID).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.UserID,
		&t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "ticket not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not fetch ticket",
		})
		return
	}

	nextOK := false
	if allowedForThisStatus, ok := allowedNext[t.Status]; ok {
		nextOK = allowedForThisStatus[input.Status]
	}

	if !nextOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid status transition",
		})
		return
	}

	updateQuery := "UPDATE tickets SET status = ? WHERE id = ? AND user_id = ?"
	_, err = database.DB.Exec(updateQuery, input.Status, ticketID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not update ticket status",
		})
		return
	}

	t.Status = input.Status

	c.JSON(http.StatusOK, gin.H{
		"message": "ticket status updated successfully",
		"ticket":  t,
	})
}

