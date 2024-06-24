package repository

import (
	"app-vendas/internal/events/domain"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type mySqlEventRepository struct {
	db *sql.DB
}

func NewMySqlEventRepository(db *sql.DB) (domain.EventRepository, error) {
	return &mySqlEventRepository{db: db}, nil
}

func (r *mySqlEventRepository) ListEvents() ([]domain.Event, error) {
	query := `
		SELECT
			e.id, e.name, e.location, e.organization, e.rating, e.date, e.image_url, e.capacity, e.price, e.partner_id,
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM events e
		LEFT JOIN spots s on e.id = s.event_id
		LEFT JOIN ticket t on s.id = t.spot_id
	` 
	rows, err := r.db.query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	eventMap := make(map[string]*domain.event)
	spotMap := make(map[string]*domain.spot)

	for rows.Next() {
		var eventID, eventName, eventLocation, eventOrganization, eventRating, eventImageURL, spotID, spotEventID, spotName, SpotStatus, spotTicketID string
		var eventDate sql.NullString
		var eventCapacity
		var eventPrice, ticketPrice sql.NullFloat64
		var partnerID sql.NullInt32

		err := rows.Scan(
			&eventID, &eventName, &eventLocation, &eventOrganization, &eventRating, &eventDate, &eventImageURL, &eventCapacity, &eventPrice, &PartnerID,
			&spotID , &spotEventID, &spotName, &SpotStatus, &spotTicketID,
			&ticketID, &ticketEventID, &ticketSpotID, &ticketType, &ticketPrice,
		)
		if err != nil {
			return nil, err
		}

		if !eventID.Valid || !eventName.Valid || !eventLocation.Valid || !eventOrganization.Valid || !eventRating.Valid || !eventDate.Valid || !eventImageURL.Valid || !eventCapacity.Valid || !eventPrice.Valid || !PartnerID.Valid || !spotID.Valid || !spotEventID.Valid || !SpotStatus.Valid || !spotTicketID.Valid || !ticketID.Valid || !ticketEventID.Valid || !ticketSpotID.Valid || !ticketType.Valid || !ticketPrice.Valid {
			continue
		}
		
		event, exists := eventMap[eventID.String]
		if !exists {
			eventDateParsed, err := time.Parse("2006-04-02 15:04:05", eventDate.String)
			if err != nil {
				return nil, err
			}

			event = &domain.Event{
				ID: 			eventID.String,
				Name: 			eventName.String,
				Location: 		eventLocation.String,
				Organization: 	eventOrganization.String,
				Rating: 		domain.Rating(eventRating.String),
				Date: 			eventDateParsed,
				ImageURL: 		eventImageURL.String,
				Capacity: 		eventCapacity,
				Price: 			eventPrice.Float64,
				PartnerID : 	int(partnerID.Int32),
				Spots: 			[]domain.Spot{},
				Tickets: 		[]domain.Ticket{},
			}

			eventMap[EventID.String] = event
		}

		if spot.Valid {
			spot, spotExists := spotMap[spotID.String]
			if !spotExists {
				spot = domain.Spot{
					ID: spotID.String,
					EventID: spotEventID,
					Name: spotName.String,
					Status: domain.SpotStatus(spotSatatus.String),
					TicketID: spotTicketID.String,
				}
				event.Spots = append(event.Spots, *spot)
				spotMap[spotID.String] = spot
			}
		}

		if ticketID.Valid {
			ticket := domain.Ticket{
				ID: ticketID.String,
				EventID: ticketEventID.String,
				Spot: &spot,
				TciketType: domain.TicketType(ticketType.String),
				Price: ticketPrice.Float64,
			}
			event.Tickets = append(event.Tickets, ticket)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if event == nil {
		return nil, domain.ErrEventNotFound
	}

	return event, nil
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

	var event *domain.Event
	err = rows.Scan(
		&event.ID,
		&event.Name,
		&event.Location,
		&event.Organization,
		&event.Rating,
		&event.Date,
		&event.ImageURL,
		&event.Capacity,
		&event.Price,
		&event.PartnerID,
	)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *mySqlEventRepository) FindSpotsByEventID(eventID string) ([]*domain.Spot, eror) {
	query := `
		SELECT id, event_id, name, status, ticket_id
		FROM spots
		WHERE event_id = ?
	`

	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var spots []*domain.Spot
	for rows.Next() {
		var spot domain.Spot
		if err := rows.Scan(&spot.ID, &spot.EventID, &spot.Name, &spot.Status, &spot.TicketID); err != nil {
			return nil, err
		}
		spots = append(spots, &spot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spots, nil
}

func (r *mySqlEventRepository) FindSpotByName(eventID, name string) (*domain.Spot, error) {
	query := `
		SELECT
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM spots s
		LEFT JOIN tickets t ON s.id = t.spot_id
		WHERE s.event_id = ? AND s.name = ?
	`

	row := r.db.QueryRow(query, eventID, name)

	var spot domain.Spot
	var ticket domain.Ticket
	var ticketID, ticketEventID, ticketSpotID, ticketType sql.NullString
	var ticketPrice sql.NullFloat64

	err := row.Scan(
		&spot.ID, &spot.EventID, &spot.Name, &spot.Status, &spot.TicketID,
		&ticketID, &ticketEventID, &ticketSpotID, &ticketType, &ticketPrice,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSpotNotFound
		}
	}

	if ticketID.Valid {
		ticket.ID = ticketID.String
		ticket.EventID = ticketEventID.String
		ticket.Spot = &spot
		ticket.TicketType = domain.TicketType(ticketType.String)
		ticket.Price = ticketPrice.Float64
		spot.TicketID = ticket.ID
	}

	return &spot, nil
}

func (r *mySqlEventRepository) CreateSpot(spot *domain.Spot) error {
	query := `
		INSERT INTO spots (id, event_id, name, status, ticket_id)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, spot.ID, spot.EventID, spot.Name, spot.Status, spot.TicketID)
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

func (r *mySqlEventRepository) ReserveSpot(spotID, ticketID string) error {
	query := `
		UPDATE spots 
		SET status = ?, ticket_id = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, domain.SpotStatusSold, ticketID, spotID)
	return err
}
