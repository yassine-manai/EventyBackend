package third_party

import (
	"eventy/pkg/db"
	"eventy/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func BookEventHandler(c *gin.Context) {
	var req models.BookEventRequest

	//log.Debug().Interface("BookEventHandler API request ", req).Send()
	// Bind JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Debug().Interface("BookEventHandler API request ", &req).Send()

	// Call the booking function
	rowsAffected, err := db.BookEvent(c.Request.Context(), req.EventID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event, err := db.GetEventByID(c.Request.Context(), req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// increment the price from user balance
	_, err = db.TopDownBalance(c.Request.Context(), req.UserID, event.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":       "Event booked successfully",
		"rows_affected": rowsAffected,
	})
}
