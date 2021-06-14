package translator

import (
	"encoding/json"
	"log"
)

type Message struct {
	msg string
}

func ToJSON(messageToSend string) string {
	message := Message{messageToSend}

	messageEncoded, err := json.Marshal(message)

	if err != nil {
		log.Fatalf("Message cannot be encoding")
	}
	return string(messageEncoded)
}

func translateMessage(messageToTranslate []byte) (error, Message) {
	var msg Message

	err := json.Unmarshal(messageToTranslate, &msg)

	if err != nil {
		log.Fatalf("Message cannot be decoding")
	}

	return err, msg
}

func SendMessage(messageToSend string) string {
	messageJson := ToJSON(messageToSend)
	return messageJson
}

func ReceiveMessage(messageToReceive []byte) interface{} {
	_, messageToString := translateMessage(messageToReceive)
	return messageToString
}