package chatroom

import (
	"sync"

	"github.com/benharmonics/personal-site-backend/database/models"
	"github.com/benharmonics/personal-site-backend/logging"
)

type messageHistory struct {
	mu       sync.Mutex
	messages [][]byte
}

func (hist *messageHistory) push(message []byte) {
	hist.mu.Lock()
	defer hist.mu.Unlock()
	hist.messages = append(hist.messages, message)
}

func (hist *messageHistory) saveToDB(roomName string) error {
	if database == nil {
		logging.Debug("database unavailable: not set")
		return nil
	}
	msgs, err := hist.toDBModels(roomName)
	if err != nil {
		return err
	}
	if err = database.InsertChatroomMessages(msgs); err != nil {
		return err
	}
	hist.mu.Lock()
	defer hist.mu.Unlock()
	hist.messages = [][]byte{} // Empty cached messages
	return nil
}

func (hist *messageHistory) toDBModels(roomName string) ([]models.ChatroomMessage, error) {
	var ret []models.ChatroomMessage
	hist.mu.Lock()
	defer hist.mu.Unlock()
	for _, msg := range hist.messages {
		msg, err := models.NewChatroomMessage(msg, roomName)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *msg)
	}
	return ret, nil
}
