package chatroom

var roomManager = make(map[string]*chatroom)

type chatroom struct {
	clients        map[*client]bool
	register       chan *client
	unregister     chan *client
	broadcast      chan []byte
	messageHistory [][]byte
}

func newChatroom(name string) *chatroom {
	if room, ok := roomManager[name]; ok {
		return room
	}
	room := &chatroom{
		clients:        make(map[*client]bool),
		register:       make(chan *client),
		unregister:     make(chan *client),
		broadcast:      make(chan []byte),
		messageHistory: [][]byte{},
	}
	roomManager[name] = room
	return room
}

func (room *chatroom) run() {
	for {
		select {
		case client := <-room.register:
			room.clients[client] = true
			for _, message := range room.messageHistory {
				client.send <- message
			}
		case client := <-room.unregister:
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.send)
			}
		case message := <-room.broadcast:
			room.messageHistory = append(room.messageHistory, message)
			for client := range room.clients {
				client.send <- message
			}
		}
	}
}
