package platform

import "time"

type Message struct {
	ID   string
	From string
	Text string
	Time time.Time
}

type Platform interface {
	SendMessage(to, text string) error
	ReadMessages(from string, limit int) ([]Message, error)
}
