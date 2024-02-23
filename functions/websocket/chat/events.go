package main

type ChatEventType string

const (
	UserSendMessage           ChatEventType = "USER:SEND_MESSAGE"
	UserUpdateMessageStatus   ChatEventType = "USER:UPDATE_MESSAGE_STATUS"
	ServerSendMessage         ChatEventType = "SERVER:SEND_MESSAGE"
	ServerUpdateMessageStatus ChatEventType = "SERVER:UPDATE_MESSAGE_STATUS"
)

type ChatEvent struct {
	Type    ChatEventType `json:"type"`
	Payload any           `json:"payload"`
}

type UserSendMessagePayload struct {
	Content        string `json:"content"`
	ConversationID string `json:"conversationId"`
}

type MessageStatus string

type UserUpdateMessageStatusPayload struct {
	Content        string `json:"content"`
	ConversationID string `json:"conversationId"`
	MessageID      string `json:"messageId"`
	Status         string
}
