package client

func sendMessage(messageToSend string) string {
	messageJson := ToJSON(messageToSend)
	return messageJson
}

func receiveMessage(messageToReceive []byte) interface{} {
	_, messageToString := translateMessage(messageToReceive)
	return messageToString
}

func AddingValue(message string) string {
	return "AAAAAAAAA.... " + message
}