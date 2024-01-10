package suggest

import (
	"fmt"
	"time"
)

const (
	DefaultTimeFormat = "02/01/2006-15:04:05"
)

type Message struct {
	Sender    string
	Receiver  string
	Content   string
	Timestamp int64 // Unix timestamp
}

func (m Message) String() string {
	return fmt.Sprintf("[From: %s, To: %s, At: %d, Content: %s]", m.Sender, m.Receiver, m.Timestamp, m.Content)
}

func NewMessage(Sender string, Receiver string, content string) *Message {
	return &Message{Sender, Receiver, content, time.Now().Unix()}
}
