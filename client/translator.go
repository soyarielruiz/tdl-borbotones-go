package client

import (
	"encoding/json"
	"log"
)

type Message struct {
	msg string
}

func ToJSON(messageToSend string) string {
	message := Message{messageToSend}

	messageEncoded, error := json.Marshal(message)

	if error != nil {
		log.Fatalf("Message cannot be encoding")
	}
	return string(messageEncoded)
}

func translateMessage(messageToTranslate []byte) (error, Message) {
	var msg Message

	error := json.Unmarshal(messageToTranslate, &msg)

	if error != nil {
		log.Fatalf("Message cannot be decoding")
	}

	return error, msg
}