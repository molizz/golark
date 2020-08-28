package messager

import (
	"github.com/molizz/golark/message/card"
)

type ExampleMessage struct {
	title, desc, url string
}

func (t *ExampleMessage) Kind() int {
	return MessageKindExample
}

func (t *ExampleMessage) Type() string {
	return MessageTypeDefault
}

func (t *ExampleMessage) Card() *card.Card {
	return card.NewSimpleCard(t.title, t.desc, t.url)
}

func NewExampleMessage(title, desc, url string) *ExampleMessage {
	return &ExampleMessage{
		title: title,
		desc:  desc,
		url:   url,
	}
}
