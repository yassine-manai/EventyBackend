package third_party

import (
	"context"
	"eventy/pkg/db"
	"eventy/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// user login
func Login(c *gin.Context) {
	var loginDetails models.Login
	ctx := context.Background()

	// Bind incoming JSON to loginDetails
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if user exists
	user, err := db.GetUserByEmail(ctx, loginDetails.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Validate password (later, you should hash passwords)
	if user.Password != loginDetails.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Return user details
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user_id": user.UserID,
		"email":   user.Email,
		"name":    user.Name,
		"token":   "Bearer token",
	})
}

// Register new user
func Register(c *gin.Context) {
	log.Debug().Msg("Register API request")

	var registerDetail models.User
	ctx := context.Background()
	registerDetail.Is_guest = true

	if err := c.ShouldBindJSON(&registerDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := db.GetUserByEmail(ctx, registerDetail.Email)
	if err != nil {

		// Check if user exists
		err := db.AddUser(ctx, &registerDetail)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Errror creating user"})
			return
		}

		return
	}
	if user.Email == registerDetail.Email {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Register successful"})

}

// Update user profile
func UpdateProfile(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Query("user_id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Err(err).Str("UserID", idStr).Msg("Invalid User ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid User ID",
			"code":    -400,
		})
		return
	}

	var updates models.User
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
			"code":    -400,
		})
		return
	}

	rowsAffected, err := db.UpdateUser(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Int("UserID", id).Msg("No user found with the given ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No user found with the given ID",
			"code":    -404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
		"code":    200,
	})
}

// Get user profile by ID
func GetUserProfile(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Query("user_id")
	user := []models.User{}

	id, _ := strconv.Atoi(idStr)
	if idStr != "" {
		log.Debug().Int("UserID", id).Msg("Get User by ID API mobile request")
		user, err := db.GetUserByID(ctx, id)
		if err != nil {
			log.Warn().Err(err).Str("UserID", idStr).Msg("Error retrieving User ID")
			c.JSON(http.StatusOK, []models.User{})
			return
		}

		c.JSON(http.StatusOK, user)
		return
	}

	c.JSON(http.StatusOK, user)
}
