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

// GetCategories godoc
//
//	@Summary		Get all categories
//	@Description	Get a list of all categories
//	@Tags			Backoffice - Categories
//	@Produce		json
//	@Param			category_id	query	string	false	"Category ID"
//	@Success		200			{array}	models.Category	"List of Categories"
//	@Router			/get_categories [get]
func GetCategories(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Query("category_id")

	id, _ := strconv.Atoi(idStr)

	if idStr != "" {
		log.Debug().Int("CategoryID", id).Msg("Get Category by ID API request")
		category, err := db.GetCategoryByID(ctx, id)
		if err != nil {
			log.Warn().Err(err).Str("CategoryID", idStr).Msg("Error retrieving Category ID")
			c.JSON(http.StatusOK, []models.Category{})
			return
		}

		c.JSON(http.StatusOK, category)
		return
	}

	categories, err := db.GetAllCategories(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all categories")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if len(categories) == 0 {
		log.Debug().Int("Category List", len(categories)).Msg("No data found")
		c.JSON(http.StatusOK, []models.Category{})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// AddCategory godoc
//
//	@Summary		Add a new category
//	@Description	Add a new category to the database
//	@Tags			Backoffice - Categories
//	@Accept			json
//	@Produce		json
//	@Param			category	body	models.Category	true	"Category data"
//	@Router			/add_category [post]
func AddCategory(c *gin.Context) {
	ctx := context.Background()
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
			"code":    -400,
		})
		return
	}

	err := db.AddCategory(ctx, &category)
	if err != nil {
		log.Err(err).Msg("Error adding category")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to add category",
			"code":    -500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category added successfully",
		"code":    200,
	})
}

// UpdateCategory godoc
//
//	@Summary		Update a category
//	@Description	Update an existing category in the database
//	@Tags			Backoffice - Categories
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path	int				true	"Category ID"
//	@Param			category	body	models.Category	true	"Updated category data"
//	@Router			/update_category/{category_id} [put]
func UpdateCategory(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Param("category_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Err(err).Str("CategoryID", idStr).Msg("Invalid Category ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Category ID",
			"code":    -400,
		})
		return
	}

	var updates models.Category
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
			"code":    -400,
		})
		return
	}

	rowsAffected, err := db.UpdateCategory(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating category")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update category",
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Int("CategoryID", id).Msg("No category found with the given ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No category found with the given ID",
			"code":    -404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category updated successfully",
		"code":    200,
	})
}

// DeleteCategory godoc
//
//	@Summary		Delete a category
//	@Description	Delete a category from the database
//	@Tags			Backoffice - Categories
//	@Produce		json
//	@Param			category_id	path	int	true	"Category ID"
//	@Router			/delete_category/{category_id} [delete]
func DeleteCategory(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Param("category_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Err(err).Str("CategoryID", idStr).Msg("Invalid Category ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Category ID",
			"code":    -400,
		})
		return
	}

	rowsAffected, err := db.DeleteCategory(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting category")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete category",
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Int("CategoryID", id).Msg("No category found with the given ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No category found with the given ID",
			"code":    -404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted successfully",
		"code":    200,
	})
}
