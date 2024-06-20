package repository

import (
	"app-vendas/internal/events/domain"
	"database/sql"
)

type mySqlEventRepository struct {
	db *sql.DB
}

// func NewMySqlEventRepository(db *sql.DB) domain.EventRepository, error {
// 	return &mySqlEventRepository{db: db}, nil
// }

func (r *mySqlEventRepository) CreateSpot(spot *domain.Spot) error {
	query := `
		INSERT INTO spots (id, event_id, name, status, ticket_id)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, spot.ID, spot.EventID, spot.Name, spot.Status, spot.TicketID)
	return err
}

func (r *mySqlEventRepository) ReserveSpot(spotID, ticketID string) error {
	query := `
		UPDATE spots 
		SET status = ?, ticket_id = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, domain.SpotStatusSold, ticketID, spotID)
	return err
}

func (r *mySqlEventRepository) CreateTicket(ticket *domain.Ticket) error {
	query := `
		INSERT INTO tickets (id, event_id, spot_id, ticket_type, price)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, ticket.ID, ticket.EventID, ticket.Spot.ID, ticket.TicketType, ticket.Price)
	return err
}

func (r *mySqlEventRepository) FindEventByID(eventID string) (*domain.Event, error) {
	query := `
		SELECT id, name, location, organization, rating, date, image_url, capacity, price, partner_id
		FROM events
		WHERE id = ?
	`

	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

}
