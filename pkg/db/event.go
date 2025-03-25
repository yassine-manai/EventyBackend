package db

import (
	"context"
	"eventy/pkg/models"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// GetAllEvents retrieves all events from the database
func GetAllEvents(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	err := Db_GlobalVar.NewSelect().Model(&events).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all events: %w", err)
	}
	return events, nil
}

// GetEventByID retrieves a single event by its ID
func GetEventByID(ctx context.Context, id int) (*models.Event, error) {
	event := new(models.Event)
	err := Db_GlobalVar.NewSelect().Model(event).Where("event_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting event by ID %d: %w", id, err)
	}
	return event, nil
}

// AddEvent creates a new event in the database
func AddEvent(ctx context.Context, event *models.Event) error {

	_, err := Db_GlobalVar.NewInsert().Model(event).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	log.Debug().Msgf("New event added with ID: %d", event.EventID)
	return nil
}

// UpdateEvent updates an existing event in the database
func UpdateEvent(ctx context.Context, id int, updates *models.Event) (int64, error) {
	res, err := Db_GlobalVar.NewUpdate().
		Model(updates).
		Where("event_id = ?", id).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating event with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated event with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}

// UpdateEvent updates an existing event in the database
func BookEvent(ctx context.Context, id int, userID int) (int64, error) {
	log.Info().Msgf("Starting booking process for Event ID: %d, User ID: %d", id, userID)

	// Fetch the current list of users for the event
	var event models.Event
	err := Db_GlobalVar.NewSelect().
		Model(&event).
		Where("event_id = ?", id).
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("Error fetching event with ID %d", id)
		return 0, fmt.Errorf("error fetching event with ID %d: %w", id, err)
	}
	log.Debug().Msgf("Fetched event: %+v", event.EventID)

	// Check if the user is already booked or if the event is full
	for _, uid := range event.UserID {
		if uid == userID {
			log.Warn().Msgf("User %d is already booked for event %d", userID, id)
			return 0, fmt.Errorf("user %d already booked for this event", userID)
		}
	}

	log.Debug().Msgf("Event capacity: %d, Current users: %d", event.Capacity, len(event.UserID))
	if len(event.UserID) >= event.Capacity {
		log.Warn().Msgf("Event ID %d is full. Capacity: %d", id, event.Capacity)
		return 0, fmt.Errorf("event is full")
	}

	// Add new user ID to the event
	event.UserID = append(event.UserID, userID)
	log.Info().Msgf("Added User ID %d to Event ID %d", userID, id)

	// Update the event with the new user list
	res, err := Db_GlobalVar.NewUpdate().
		Model(&event).
		Where("event_id = ?", id).
		Set("user_id = ?", pgdialect.Array(event.UserID)).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("Error updating event with ID %d", id)
		return 0, fmt.Errorf("error updating event with ID %d: %w", id, err)
	}
	log.Info().Msgf("Updated Event ID %d with new user list", id)

	// Fetch the user
	var user models.User
	err = Db_GlobalVar.NewSelect().
		Model(&user).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("Error fetching user with ID %d", userID)
		return 0, fmt.Errorf("error fetching user with ID %d: %w", userID, err)
	}
	log.Debug().Msgf("Fetched user: %+v", user)

	// Add event ID to the user's EventID list
	user.EventID = append(user.EventID, id)
	log.Info().Msgf("Added Event ID %d to User ID %d", id, userID)

	// Update the user with the new event list
	_, err = Db_GlobalVar.NewUpdate().
		Model(&user).
		Where("user_id = ?", userID).
		//Set("user_id = ?", pgdialect.Array(event.UserID)).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("Error updating user with ID %d", userID)
		return 0, fmt.Errorf("error updating user with ID %d: %w", userID, err)
	}
	log.Info().Msgf("Updated User ID %d with new event list", userID)

	// Log and return the result
	rowsAffected, _ := res.RowsAffected()
	log.Info().Msgf("Successfully booked User ID %d for Event ID %d, Rows Affected: %d", userID, id, rowsAffected)

	return rowsAffected, nil
}

// DeleteEvent removes an event from the database by its ID
func DeleteEvent(ctx context.Context, id int) (int64, error) {
	res, err := Db_GlobalVar.NewDelete().Model(&models.Event{}).Where("event_id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting event with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted event with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}
