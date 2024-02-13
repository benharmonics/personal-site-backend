package models

import (
	"fmt"
	"strings"
	"time"
)

type ChatroomMessage struct {
	RoomName  string    `bson:"roomName"`
	Author    string    `bson:"author"`
	Message   string    `bson:"message"`
	Timestamp time.Time `bson:"ts"`
}

func NewChatroomMessage(rawMessage []byte, roomName string) (*ChatroomMessage, error) {
	data := strings.Split(string(rawMessage), " ")
	if len(data) < 3 {
		return nil, fmt.Errorf("malformed message: %s", string(rawMessage))
	}
	iLast := len(data) - 1
	timestamp, err := time.Parse(time.RFC3339, data[iLast])
	if err != nil {
		return nil, err
	}
	msg := &ChatroomMessage{
		RoomName:  roomName,
		Author:    strings.Trim(data[0], ":"),
		Message:   strings.Join(data[1:iLast], " "),
		Timestamp: timestamp,
	}
	return msg, nil
}
