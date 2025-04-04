package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/devgymbr/kafka"
	"github.com/google/uuid"
)

func Publish(conn net.Conn, body, topic string) error {
	commandBody := kafka.Message{
		Headers: map[string]string{
			"id": uuid.NewString(),
		},
		Body: body,
	}

	bodyRaw, err := json.Marshal(commandBody)
	if err != nil {
		return err
	}

	cmd := kafka.Command{
		Type:  kafka.TypePublish,
		Topic: topic,
		Body:  string(bodyRaw),
	}

	raw, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(conn, string(raw))

	return err
}

func Consume(conn net.Conn, topic, consumerName string) (chan kafka.Message, error) {
	cmd := kafka.Command{
		Type:         kafka.TypeConsume,
		Topic:        topic,
		ConsumerName: consumerName,
	}
	messages := make(chan kafka.Message)
	raw, err := json.Marshal(cmd)
	if err != nil {
		return messages, err
	}

	if _, err = fmt.Fprintln(conn, string(raw)); err != nil {
		return messages, err
	}

	reader := bufio.NewReader(conn)
	go func() {
		defer close(messages)
		for {
			reply, _, err := reader.ReadLine()
			if err == io.EOF {
				return
			}

			if err != nil {
				println(err)
				return
			}

			var response kafka.Response
			if err := json.Unmarshal(reply, &response); err != nil {
				continue
			}
			var message kafka.Message
			if err := json.Unmarshal([]byte(response.Body), &message); err != nil {
				continue
			}

			messages <- message
		}
	}()

	return messages, nil
}
