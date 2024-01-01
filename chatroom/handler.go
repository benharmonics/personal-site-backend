package chatroom

import (
	"net/http"

	"github.com/benharmonics/personal-site-backend/logging"
	"github.com/benharmonics/personal-site-backend/utils/web"
)

func ServeChatroom(name string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("Upgrader failed to upgrade:", err)
		web.HTTPError(w, http.StatusInternalServerError)
		return
	}
	room := newChatroom(name)
	go room.run()

	client := newClient(room, conn)
	client.room.register <- client
	go client.readPump()
	go client.writePump()
}
