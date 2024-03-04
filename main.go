package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/benharmonics/personal-site-backend/api"
	cfg "github.com/benharmonics/personal-site-backend/config"
	db "github.com/benharmonics/personal-site-backend/database"
	"github.com/benharmonics/personal-site-backend/logging"
	"github.com/benharmonics/personal-site-backend/utils"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logging.Warn("Failed to load dotenv file:", err)
	}
	utils.Must(cfg.ValidateConfig())
}

func main() {
	logging.SetLogLevel(logging.LogLevelDebug)
	logging.SetColor(true)
	logging.SetTime(true)
	logging.Info("Starting")

	dbConf := cfg.NewMongoConfig()
	database, err := db.NewDatabase(
		db.WithEncryptedConnection(),
		db.WithHost(dbConf.Host),
		db.WithCredentials(dbConf.Username, dbConf.Password),
		db.WithoutPort(),
	)
	utils.Must(err)
	defer database.Disconnect()

	srv := api.NewServer(database)
	appConf := cfg.NewAppConfig()
	addr := fmt.Sprintf("%s:%d", appConf.Host, appConf.Port)
	logging.Info("Listening on", addr)
	utils.Must(http.ListenAndServe(addr, srv))
}
