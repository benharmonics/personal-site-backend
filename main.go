package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/benharmonics/personal-site-backend/api"
	"github.com/benharmonics/personal-site-backend/config"
	"github.com/benharmonics/personal-site-backend/logging"
	"github.com/benharmonics/personal-site-backend/utils"
)

func init() {
	logging.SetColor(true)
	logging.SetDebug(true)
	if err := godotenv.Load(); err != nil {
		logging.Warning("Failed to load dotenv file:", err)
	}
	utils.Must(config.ValidateConfig())
}

func main() {
	srv := api.NewServer()
	defer srv.DisconnectFromDatabase()
	appConf := config.NewAppConfig()
	addr := fmt.Sprintf("%s:%d", appConf.Host, appConf.Port)
	logging.Info("Listening on", addr)
	utils.Must(http.ListenAndServe(addr, srv))
	http.ListenAndServe(addr, srv)
}
