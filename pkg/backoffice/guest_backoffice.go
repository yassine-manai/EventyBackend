package backoffice

import (
	"context"
	"eventy/pkg/db"
	"eventy/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetGuests godoc
//
//	@Summary		Get all guests
//	@Description	Get a list of all guests
//	@Tags			Backoffice - Guests
//	@Produce		json
//	@Success		200		{array}	models.User	"List of Guests"
//	@Router			/get_guests [get]
func GetGuests(c *gin.Context) {
	ctx := context.Background()

	guest, err := db.GetAllGuests(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all guests")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if len(guest) == 0 {
		log.Debug().Int("Guest List", len(guest)).Msg("No data found")
		c.JSON(http.StatusOK, []models.User{})
		return
	}

	c.JSON(http.StatusOK, guest)
}

// DeclineGuets godoc
//
//	@Summary		Decmine a user
//	@Description	Decmine a user from the database
//	@Tags			Backoffice - Guets
//	@Produce		json
//	@Param			user_id	path	int	true	"Guest ID"
//	@Router			/accept_guest/{user_id} [post]
func AcceptGuest(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Param("user_id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Warn().Err(err).Str("GuestID", idStr).Msg("Invalid Guest ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Guest ID",
			"code":    -400,
		})
		return
	}

	rows, err := db.AcceptGuest(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error adding user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to add user",
			"code":    -500,
		})
		return
	}

	if rows == 0 {
		log.Warn().Int("GuestID", id).Msg("No Guest found with the given ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No Guest found with the given ID",
			"code":    -404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Guest Accepted successfully",
		"code":    200,
	})
}

// DeclineGuets godoc
//
//	@Summary		Decmine a user
//	@Description	Decmine a user from the database
//	@Tags			Backoffice - Guets
//	@Produce		json
//	@Param			user_id	path	int	true	"Guest ID"
//	@Router			/decline_user/{user_id} [delete]
func DeclineGuest(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Param("user_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Err(err).Str("GuestID", idStr).Msg("Invalid Guest ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Guest ID",
			"code":    -400,
		})
		return
	}

	rowsAffected, err := db.DeclineGuest(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to decline guest",
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Int("GuestID", id).Msg("No Guest found with the given ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No Guest found with the given ID",
			"code":    -404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Guest declined successfully",
		"code":    200,
	})
}
