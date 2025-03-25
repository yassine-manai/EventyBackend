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

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Get a list of all users
//	@Tags			Backoffice - Users
//	@Produce		json
//	@Param			user_id	query	string	false	"UserID"
//	@Success		200		{array}	models.User	"List of Users"
//	@Router			/get_users [get]
func GetUsers(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Query("user_id")

	id, _ := strconv.Atoi(idStr)

	if idStr != "" {
		log.Debug().Int("UserID", id).Msg("Get User by ID API request")
		user, err := db.GetUserByID(ctx, id)
		if err != nil {
			log.Warn().Err(err).Str("UserID", idStr).Msg("Error retrieving User ID")
			c.JSON(http.StatusOK, []models.User{})
			return
		}

		c.JSON(http.StatusOK, user)
		return
	}

	users, err := db.GetAllUsers(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if len(users) == 0 {
		log.Debug().Int("User List", len(users)).Msg("No data found")
		c.JSON(http.StatusOK, []models.User{})
		return
	}

	c.JSON(http.StatusOK, users)
}

// AddUser godoc
//
//	@Summary		Add a new user
//	@Description	Add a new user to the database
//	@Tags			Backoffice - Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	models.User	true	"User data"
//	@Router			/add_user [post]
func AddUser(c *gin.Context) {
	ctx := context.Background()
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
			"code":    -400,
		})
		return
	}

	err := db.AddUser(ctx, &user)
	if err != nil {
		log.Err(err).Msg("Error adding user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to add user",
			"code":    -500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User added successfully",
		"code":    200,
	})
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update an existing user in the database
//	@Tags			Backoffice - Users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path	int		true	"User ID"
//	@Param			user	body	models.User	true	"Updated user data"
//	@Router			/update_user/{user_id} [put]
func UpdateUser(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Param("user_id")
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

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user from the database
//	@Tags			Backoffice - Users
//	@Produce		json
//	@Param			user_id	path	int	true	"User ID"
//	@Router			/delete_user/{user_id} [delete]
func DeleteUser(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Param("user_id")
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

	rowsAffected, err := db.DeleteUser(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete user",
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
		"message": "User deleted successfully",
		"code":    200,
	})
}

func TopupBalance(c *gin.Context) {
	ctx := context.Background()
	//var user models.User
	idStr := c.Query("user_id")
	balanceStr := c.Query("balance")

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

	balance, err := strconv.Atoi(balanceStr)
	if err != nil {
		log.Warn().Err(err).Str("UserID", idStr).Msg("Invalid User ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid User ID",
			"code":    -400,
		})
		return
	}

	rowsAffected, err := db.TopupBalance(ctx, id, balance)
	if err != nil {
		log.Err(err).Msg("Error occured")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to topup balance",
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
		"message": "User Balance Topuped successfully",
		"code":    200,
	})
}
