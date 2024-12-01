package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"

	"github.com/SerzhLimon/PaymentService/config"
	serv "github.com/SerzhLimon/PaymentService/internal/transport"
	"github.com/SerzhLimon/PaymentService/pkg/postgres"
	"github.com/SerzhLimon/PaymentService/pkg/postgres/migrations"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("Loading configuration...")
	cfg := config.LoadConfig()
	logrus.Debugf("Configuration loaded: %+v", cfg)

	logrus.Info("Initializing PostgreSQL client...")
	db, err := postgres.InitPostgresClient(cfg.Postgres)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize PostgreSQL client")
	}
	defer func() {
		logrus.Info("Closing PostgreSQL connection...")
		db.Close()
		logrus.Info("PostgreSQL connection closed")
	}()
	logrus.Info("PostgreSQL client initialized successfully")

	logrus.Info("Running migrations...")
	err = migrations.Up(db)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to apply migrations")
	}
	logrus.Info("Migrations applied successfully")

	logrus.Info("Initializing server...")
	server := serv.NewServer(db)
	routes := serv.ApiHandleFunctions{
		Server: *server,
	}

	logrus.Info("Setting up router...")
	router := serv.NewRouter(routes)

	logrus.Infof("Starting server on port %s...", ":8080")
	if err := router.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
	logrus.Info("Server started successfully")
}
