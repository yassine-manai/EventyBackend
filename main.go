package main

import (
	"context"
	"database/sql"
	"eventy/config"
	_ "eventy/docs"
	"eventy/functions"
	"eventy/pkg/db"
	"eventy/pkg/models"
	"eventy/routes"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// @title						Eventy
// @version						1.00.
// @securityDefinitions.apikey	BearerAuth3rdParty
// @in							header
// @name						Authorization
// @description				Authorization token for third-party section (Ensure the token is in this format: Bearer token)
// @securityDefinitions.apikey	BearerAuthBackOffice
// @in							header
// @name						Authorization
// @description				Authorization token for back-office section (Ensure the token is in this format: Bearer token)
func main() {
	config.InitLogger()
	ctx := context.Background()

	log.Info().Msg("------------------------------ # STARTING APPLICATION # ------------------------------")

	//docs.SwaggerInfo.BasePath = config.Configvar.App.SwaggerBasePath

	log.Info().Msgf("Server running on %s:%d ", config.Configvar.Server.Host, config.Configvar.Server.Port)
	log.Info().Msgf("Database connecting to %s:%d", config.Configvar.Database.Host, config.Configvar.Database.Port)

	//docs.SwaggerInfo.BasePath = fmt.Sprintf("/%s", config.Configvar.Database.Prefix)
	//log.Debug().Msg(docs.SwaggerInfo.BasePath)

	var dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", config.Configvar.Database.User, config.Configvar.Database.Password, config.Configvar.Database.Host, config.Configvar.Database.Port, config.Configvar.Database.Name, config.Configvar.Database.SSLMode)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db.Db_GlobalVar = bun.NewDB(sqldb, pgdialect.New())

	log.Debug().Msg("------------------------------- # CONNECT TO DATABASE # ------------------------------")

	// Ping the database to check connectivity
	if err := db.Db_GlobalVar.Ping(); err != nil {
		log.Warn().Str("Error", err.Error()).Msgf("Error connecting to database %v", err)
	} else {
		log.Info().Str("Database ", config.Configvar.Database.Name).Msg("Successfully Connected to the Database.")
	}

	models := []interface{}{
		&models.User{},
		&models.Event{},
		&models.Category{},
	}

	if err := functions.CreateTables(ctx, db.Db_GlobalVar, models); err != nil {
		log.Error().Err(err).Msg("Failed to create tables")
	} else {
		log.Info().Msg("Tables Created successfully.")

	}

	// Router Setup
	r := routes.SetupRouter()

	// Server Setup
	var host = fmt.Sprintf("%s:%d", config.Configvar.Server.Host, config.Configvar.Server.Port)

	log.Info().Msgf("-------------------------------- # Server running on %s # ------------------------------", host)

	if err := r.Run(host); err != nil {
		log.Err(err).Msgf("Failed to run server: %v", err)
	}

	log.Debug().Msgf("-------------------------------- # END PROGRAM # ------------------------------")
}
