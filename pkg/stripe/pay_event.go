package stripe

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentintent"
)

// PayEvent handles the payment logic
func PayEvent(c *gin.Context) {
	// Get the event_id, price, and user_id from the request body
	var req struct {
		EventID string `json:"event_id" binding:"required"`
		Price   string `json:"price" binding:"required"`
		UserID  string `json:"user_id" binding:"required"`
	}

	// Bind request data
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate input
	if req.EventID == "" || req.Price == "" || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event_id, price, and user_id are required"})
		return
	}

	// Convert price to an integer (assume it's in cents for Stripe)
	priceInt, err := strconv.Atoi(req.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}

	// Create a PaymentIntent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(priceInt)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	// Optional: Attach user_id as metadata (you can query this later in Stripe dashboard)
	params.AddMetadata("user_id", req.UserID)

	// Create PaymentIntent
	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the client secret so the frontend can complete the payment
	c.JSON(http.StatusOK, gin.H{
		"client_secret": pi.ClientSecret,
		"event_id":      req.EventID,
		"user_id":       req.UserID, // Optionally return the user_id in the response
	})
}
