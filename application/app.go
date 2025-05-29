package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/frencius/loan-service/configuration"
	_ "github.com/lib/pq"
)

type App struct {
	Config *configuration.Configuration
	DB     *sql.DB
}

func SetupApp(ctx context.Context) (*App, error) {
	app := &App{}

	// setup config
	config, err := configuration.LoadConfig()
	if err != nil {
		log.Println("failed to load config")
		return nil, err
	}
	app.Config = &config

	// setup database
	db, err := CreateDBConnection(app.Config.Database)
	if err != nil {
		log.Println("failed to create DB connection", err)
		return nil, err
	}
	app.DB = db

	return app, nil
}

func (app *App) Close() error {
	app.DB.Close()

	return nil
}

func CreateDBConnection(dbConf configuration.Database) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password='%s' dbname=%s search_path=%s sslmode=%s",
		dbConf.Host,
		dbConf.Port,
		dbConf.Username,
		dbConf.Password,
		dbConf.Name,
		dbConf.Schema,
		dbConf.SSLMode,
	)

	return sql.Open(dbConf.Driver, dsn)
}
