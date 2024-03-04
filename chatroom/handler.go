package chatroom

import (
	"net/http"

	db "github.com/benharmonics/personal-site-backend/database"
	"github.com/benharmonics/personal-site-backend/logging"
	"github.com/benharmonics/personal-site-backend/utils/web"
)

var database *db.Database

func SetDatabase(newDb *db.Database) { database = newDb }

func ServeChatroom(name string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("Upgrader failed to upgrade:", err)
		web.HTTPError(w, http.StatusInternalServerError)
		return
	}
	room := newChatroom(name)

	client := newClient(room, conn)
	client.room.register <- client
	go client.readPump()
	go client.writePump()
}
