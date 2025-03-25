package routes

import (
	"eventy/routes/backoffice_routes"
	"eventy/routes/third_party_routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()
	//router.Use(middleware.AuditMiddleware())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
	}))

	log.Debug().Msg("--------------------------  START ROUTING  ----------------------")

	backoffice_routes.Backoffice_Routes(router)
	third_party_routes.ThirdParty_Routes(router)

	//authorizedBackOffice := router.Group("/")
	//authorizedBackOffice.Use(middleware.TokenMiddlewareBackOffice())

	//backoffice_routes.BackOfficeToken(router)                // Token Generator FOR BACKOFFICE ------------------
	//backoffice_routes.BackOfficeRouter(authorizedBackOffice) // BACKOFFICE ROUTES --------------------------------

	//backoffice_routes.ExportBackoffice(authorizedBackOffice) // BACKOFFICE EXPORT ROUTES ---------------------------

	//log.Debug().Msg("--------------------------  END ROUTING  ---------------------- ")

	// Swagger Endpoint
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/docs/index.html")
	})
	router.GET("/docs/*.any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
