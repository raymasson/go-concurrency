package main

import (
	"fmt"
)

func main() {
	msgCh := make(chan message, 1)
	errCh := make(chan failedMessage, 1)

	msg := message{
		To:      []string{"frodo@underhill.me"},
		From:    "gandalf@whitecouncil.org",
		Content: "Keep it secret, keep it safe",
	}

	failedMessage := failedMessage{
		ErrorMessage:    "Message intercepted by black rider",
		OriginalMessage: message{},
	}

	msgCh <- msg
	errCh <- failedMessage

	select {
	case receivedMsg := <-msgCh:
		fmt.Println(receivedMsg)
	case receivedError := <-errCh:
		fmt.Println(receivedError)
	default:
		fmt.Println("No messages received")
	}
}

type message struct {
	To      []string
	From    string
	Content string
}

type failedMessage struct {
	ErrorMessage    string
	OriginalMessage message
}
