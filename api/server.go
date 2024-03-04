package api

import (
	"net/http"
	"time"

	"github.com/benharmonics/personal-site-backend/chatroom"
	db "github.com/benharmonics/personal-site-backend/database"
)

type Server struct {
	*http.ServeMux
	db        *db.Database
	startTime time.Time
}

func NewServer(database *db.Database) Server {
	// We have to set the chatroom database at some point or chat messages will never get saved
	chatroom.SetDatabase(database)
	srv := Server{
		ServeMux:  http.NewServeMux(),
		db:        database,
		startTime: time.Now(),
	}
	srv.routes()
	return srv
}

func (s *Server) Uptime() time.Duration { return time.Since(s.startTime) }
