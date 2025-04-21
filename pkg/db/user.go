package db

import (
	"context"
	"eventy/pkg/models"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

// GetAllUsers retrieves all users from the database
func GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := Db_GlobalVar.NewSelect().Model(&users).Where("is_guest = ?", false).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	// Remove duplicates from EventID for each user
	for i := range users {
		users[i].EventID = removeDuplicates(users[i].EventID)
	}

	return users, nil
}

// Helper function to remove duplicates from a slice of integers
func removeDuplicates(slice []int) []int {
	if len(slice) == 0 {
		return slice
	}

	// Use a map to track unique elements
	seen := make(map[int]bool)
	result := make([]int, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func GetAllGuests(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := Db_GlobalVar.NewSelect().Model(&users).Where("is_guest = ?", true).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}
	return users, nil
}

// GetUserByID retrieves a single user by their ID
func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := new(models.User)
	err := Db_GlobalVar.NewSelect().Model(user).
		Where("user_id = ?", id).
		Where("is_guest = ?", false).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting user by ID %d: %w", id, err)
	}

	// Remove duplicates from EventID
	user.EventID = removeDuplicates(user.EventID)

	if len(user.EventID) == 0 {
		return user, nil
	}

	// Get all events at once
	var events []models.Event
	err = Db_GlobalVar.NewSelect().Model(&events).
		Where("event_id IN (?)", bun.In(user.EventID)).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting events: %w", err)
	}

	// Create a map of valid event IDs (those that haven't ended)
	validEvents := make(map[int]bool)
	now := time.Now()
	for _, event := range events {
		endDate, err := time.Parse("2006-01-02", event.EndDate)
		if err != nil {
			continue
		}
		if endDate.After(now) {
			validEvents[event.EventID] = true
		}
	}

	// Filter the user's EventID list
	filteredEventIDs := make([]int, 0, len(user.EventID))
	for _, eventID := range user.EventID {
		if validEvents[eventID] {
			filteredEventIDs = append(filteredEventIDs, eventID)
		}
	}

	user.EventID = filteredEventIDs

	return user, nil
}
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	err := Db_GlobalVar.NewSelect().Model(user).
		Where("email = ?", email).
		Where("is_guest = ?", false).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting user by Email %s: %w", email, err)
	}
	user.EventID = removeDuplicates(user.EventID)

	return user, nil
}

// AddUser creates a new user in the database
func AddUser(ctx context.Context, user *models.User) error {
	user.Is_guest = true
	user.UserID = 13
	_, err := Db_GlobalVar.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	log.Debug().Msgf("New user added with Email: %s", user.Email)
	return nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(ctx context.Context, id int, updates *models.User) (int64, error) {
	res, err := Db_GlobalVar.NewUpdate().
		Model(updates).
		Where("user_id = ?", id).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating user with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated user with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}

func TopupBalance(ctx context.Context, id, Balance int) (int64, error) {
	var updates models.User
	res, err := Db_GlobalVar.NewUpdate().
		Model(&updates).
		Set("balance = balance + ?", Balance). // ✅ Increment the balance
		Where("user_id = ?", id).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error updating user with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated user with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}

func TopDownBalance(ctx context.Context, id, Balance int) (int64, error) {
	var updates models.User
	res, err := Db_GlobalVar.NewUpdate().
		Model(&updates).
		Set("balance = balance - ?", Balance).
		Where("user_id = ?", id).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error updating user with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated user with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}

// DeleteUser removes a user from the database by their ID
func DeleteUser(ctx context.Context, id int) (int64, error) {
	res, err := Db_GlobalVar.NewDelete().Model(&models.User{}).Where("user_id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting user with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted user with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}

func DeclineGuest(ctx context.Context, id int) (int64, error) {
	res, err := Db_GlobalVar.NewDelete().Model(&models.User{}).Where("user_id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error decline guest with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Declined user with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}

func AcceptGuest(ctx context.Context, id int) (int64, error) {
	var updates models.User
	res, err := Db_GlobalVar.NewUpdate().
		Model(&updates).
		Where("user_id = ?", id).
		Set("is_guest =?", false).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating user with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated user with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}
