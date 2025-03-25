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

// GetEvents godoc
//
//	@Summary		Get all events
//	@Description	Get a list of all events
//	@Tags			Backoffice - Events
//	@Produce		json
//	@Param			event_id	query	string		false	"Event ID"
//	@Success		200			{array}	db.Event	"List of Events"
//	@Router			/get_events [get]
func GetEvents(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Query("event_id")

	id, _ := strconv.Atoi(idStr)

	if idStr != "" {
		log.Debug().Int("EventID", id).Msg("Get Event by ID API request")
		event, err := db.GetEventByID(ctx, id)
		if err != nil {
			log.Warn().Err(err).Str("EventID", idStr).Msg("Error retrieving Event ID")
			c.JSON(http.StatusOK, []models.Event{})
			return
		}

		c.JSON(http.StatusOK, event)
		return
	}

	events, err := db.GetAllEvents(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all events")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if len(events) == 0 {
		log.Debug().Int("Event List", len(events)).Msg("No data found")
		c.JSON(http.StatusOK, []models.Event{})
		return
	}

	c.JSON(http.StatusOK, events)
}

// AddEvent godoc
//
//	@Summary		Add a new event
//	@Description	Add a new event to the database
//	@Tags			Backoffice - Events
//	@Accept			json
//	@Produce		json
//	@Param			event	body	db.Event	true	"Event data"
//	@Router			/add_event [post]
func AddEvent(c *gin.Context) {
	ctx := context.Background()
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
			"code":    -400,
		})
		return
	}

	//	eventCap, _ := strconv.Atoi(event.Capacity)
	err := db.AddEvent(ctx, &event)
	if err != nil {
		log.Err(err).Msg("Error adding event")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to add event",
			"code":    -500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event added successfully",
		"code":    200,
	})
}

// UpdateEvent godoc
//
//	@Summary		Update an event
//	@Description	Update an existing event in the database
//	@Tags			Backoffice - Events
//	@Accept			json
//	@Produce		json
//	@Param			event_id	path	int			true	"Event ID"
//	@Param			event		body	db.Event	true	"Updated event data"
//	@Router			/update_event/{event_id} [put]
func UpdateEvent(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Param("event_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Err(err).Str("EventID", idStr).Msg("Invalid Event ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Event ID",
			"code":    -400,
		})
		return
	}

	var updates models.Event
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
			"code":    -400,
		})
		return
	}

	rowsAffected, err := db.UpdateEvent(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating event")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update event",
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Int("EventID", id).Msg("No event found with the given ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No event found with the given ID",
			"code":    -404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event updated successfully",
		"code":    200,
	})
}

// DeleteEvent godoc
//
//	@Summary		Delete an event
//	@Description	Delete an event from the database
//	@Tags			Backoffice - Events
//	@Produce		json
//	@Param			event_id	path	int	true	"Event ID"
//	@Router			/delete_event/{event_id} [delete]
func DeleteEvent(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Param("event_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Err(err).Str("EventID", idStr).Msg("Invalid Event ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Event ID",
			"code":    -400,
		})
		return
	}

	rowsAffected, err := db.DeleteEvent(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting event")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete event",
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Int("EventID", id).Msg("No event found with the given ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No event found with the given ID",
			"code":    -404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event deleted successfully",
		"code":    200,
	})
}
