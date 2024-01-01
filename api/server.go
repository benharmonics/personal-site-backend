package api

import (
	"net/http"
	"time"

	cfg "github.com/benharmonics/personal-site-backend/config"
	db "github.com/benharmonics/personal-site-backend/database"
)

type Server struct {
	*http.ServeMux
	db        db.Database
	startTime time.Time
}

func NewServer() Server {
	dbConf := cfg.NewMongoConfig()
	database := db.NewDatabase(
		db.WithEncryptedConnection(),
		db.WithHost(dbConf.Host),
		db.WithCredentials(dbConf.Username, dbConf.Password),
		db.WithoutPort(),
	)
	srv := Server{
		ServeMux:  http.NewServeMux(),
		db:        database,
		startTime: time.Now(),
	}
	srv.routes()
	return srv
}

func (s *Server) DisconnectFromDatabase() { s.db.Disconnect() }
func (s *Server) Uptime() time.Duration   { return time.Since(s.startTime) }
