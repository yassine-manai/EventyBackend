package third_party_routes

import (
	"eventy/pkg/backoffice"
	"eventy/pkg/stripe"
	"eventy/pkg/third_party"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the backoffice package
func ThirdParty_Routes(router *gin.Engine) {

	// Category routes
	mobile_grp := router.Group("/mobile")
	{
		mobile_grp.GET("/get_events", backoffice.GetEvents)
		mobile_grp.GET("/get_categories", backoffice.GetCategories)
		mobile_grp.POST("/login", third_party.Login)
		mobile_grp.POST("/register", third_party.Register)
		mobile_grp.PUT("/update_profile", third_party.UpdateProfile)
		mobile_grp.GET("/get_profile", third_party.GetUserProfile)
		mobile_grp.POST("/book-event", third_party.BookEventHandler)
		mobile_grp.PUT("/topup_balance", backoffice.TopupBalance)
		mobile_grp.POST("/pay", stripe.PayEvent)

	}

}
