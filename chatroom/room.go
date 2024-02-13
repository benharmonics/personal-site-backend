package chatroom

import (
	"fmt"
	"sync"
	"time"

	"github.com/benharmonics/personal-site-backend/database"
	"github.com/benharmonics/personal-site-backend/database/models"
	"github.com/benharmonics/personal-site-backend/logging"
)

var appDatabase *database.Database

var existingChatrooms = make(map[string]*chatroom)

func SetDatabase(db *database.Database) { appDatabase = db }

type chatroom struct {
	name       string
	clients    map[*client]bool
	register   chan *client
	unregister chan *client
	broadcast  chan []byte
	history    messageHistory
}

type messageHistory struct {
	mu       sync.Mutex
	messages [][]byte
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

func (hist *messageHistory) push(message []byte) {
	hist.mu.Lock()
	defer hist.mu.Unlock()
	hist.messages = append(hist.messages, message)
}

func (room *chatroom) updateNewClient(c *client) {
	room.history.mu.Lock()
	defer room.history.mu.Unlock()
	msgs, err := appDatabase.FindChatroomMessages(room.name)
	if err != nil {
		logging.Error("Failed to get chatroom", room.name, "message history:", err)
	}
	// TODO: parse messages from database
	logging.Debug("Found", len(msgs), "messages for chatroom", room.name)
	for _, message := range room.history.messages {
		c.send <- message
	}
}

func (room *chatroom) saveToDB() error {
	if appDatabase == nil {
		return fmt.Errorf("database unavailable: no database set")
	}
	room.history.mu.Lock()
	defer room.history.mu.Unlock()
	for _, msg := range room.history.messages {
		data, err := models.NewChatroomMessage(msg, room.name)
		if err != nil {
			return err
		}
		if err := appDatabase.InsertChatroomMessage(data); err != nil {
			return err
		}
	}
	room.history.messages = [][]byte{} // Empty cached messages
	return nil
}

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
			if err := room.saveToDB(); err != nil {
				logging.Error("Chatroom", room.name, "failed to save to database:", err)
			}
		}
	}
}
