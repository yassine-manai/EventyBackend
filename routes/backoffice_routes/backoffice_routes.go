package backoffice_routes

import (
	"eventy/pkg/backoffice"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the backoffice package
func Backoffice_Routes(router *gin.Engine) {
	// User routes
	backoffice_grp := router.Group("/backoffice")
	{

		// User routes
		backoffice_grp.GET("/get_users", backoffice.GetUsers)
		backoffice_grp.POST("/add_user", backoffice.AddUser)
		backoffice_grp.PUT("/update_user/:user_id", backoffice.UpdateUser)
		backoffice_grp.DELETE("/delete_user/:user_id", backoffice.DeleteUser)

		// Category routes
		backoffice_grp.GET("/get_categories", backoffice.GetCategories)
		backoffice_grp.POST("/add_category", backoffice.AddCategory)
		backoffice_grp.PUT("/update_category/:category_id", backoffice.UpdateCategory)
		backoffice_grp.DELETE("/delete_category/:category_id", backoffice.DeleteCategory)

		// Event routes
		backoffice_grp.GET("/get_events", backoffice.GetEvents)
		backoffice_grp.POST("/add_event", backoffice.AddEvent)
		backoffice_grp.PUT("/update_event/:event_id", backoffice.UpdateEvent)
		backoffice_grp.DELETE("/delete_event/:event_id", backoffice.DeleteEvent)

		// Guests routes
		backoffice_grp.GET("/get_guests", backoffice.GetGuests)
		backoffice_grp.POST("/accept_guest/:user_id", backoffice.AcceptGuest)
		backoffice_grp.POST("/decline_guest/:user_id", backoffice.DeclineGuest)

	}

}
