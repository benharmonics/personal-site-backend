package chatroom

import (
	"fmt"
	"time"

	"github.com/benharmonics/personal-site-backend/logging"
)

var existingChatrooms = make(map[string]*chatroom)

type chatroom struct {
	name       string
	clients    map[*client]bool
	register   chan *client
	unregister chan *client
	broadcast  chan []byte
	history    messageHistory
}

func newChatroom(name string) *chatroom {
	if room, ok := existingChatrooms[name]; ok {
		return room
	}
	room := &chatroom{
		name:       name,
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan []byte),
		history:    messageHistory{messages: [][]byte{}},
	}
	existingChatrooms[name] = room
	go room.run()
	return room
}

func (room *chatroom) updateNewClient(c *client) {
	room.history.mu.Lock()
	defer room.history.mu.Unlock()
	msgs, err := database.FindChatroomMessages(room.name)
	if err != nil {
		logging.Error("Failed to get chatroom", room.name, "message history:", err)
		return
	}
	// TODO: parse messages from database
	logging.Debug("Found", len(msgs), "messages for chatroom", room.name)
	for _, msg := range msgs {
		data := fmt.Sprintf("%s: %s %s", msg.Author, msg.Message, msg.Timestamp.Format(time.RFC3339))
		logging.Debug("Sending message", data)
		c.send <- []byte(data)
	}
}

// func (room *chatroom) saveToDB() error {
// 	if database == nil {
// 		return fmt.Errorf("database unavailable: no database set")
// 	}
// 	msgs, err := room.history.toDBModels(room.name)
// 	if err != nil {
// 		return err
// 	}
// 	if err = database.InsertChatroomMessages(msgs); err != nil {
// 		return err
// 	}
// 	room.history.mu.Lock()
// 	defer room.history.mu.Unlock()
// 	room.history.messages = [][]byte{} // Empty cached messages
// 	return nil
// }

func (room *chatroom) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case wsClient := <-room.register:
			room.clients[wsClient] = true
			room.updateNewClient(wsClient)
		case wsClient := <-room.unregister:
			if _, ok := room.clients[wsClient]; ok {
				delete(room.clients, wsClient)
				close(wsClient.send)
			}
		case message := <-room.broadcast:
			room.history.push(message)
			for client := range room.clients {
				client.send <- message
			}
		case <-ticker.C:
			logging.Debug("Saving messages from chatroom", room.name, "to database")
			if err := room.history.saveToDB(room.name); err != nil {
				logging.Error("Chatroom", room.name, "failed to save to database:", err)
			}
		}
	}
}
