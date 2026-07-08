package models

// Ticket represents one row from the tickets table.
type Ticket struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      int    `json:"user_id"`
	CreatedAt   string `json:"created_at"`
}

