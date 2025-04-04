package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/devgymbr/kafka"
)

var consumers map[string]kafka.Consumer

type Config struct {
	Path    string
	Workers uint
}

func Start(conf Config, listen *net.TCPListener, done <-chan struct{}) {
	consumers = make(map[string]kafka.Consumer)
	commands := make(chan kafka.Command)
	stopCommands := make(chan bool, 1)

	go waitForCommands(listen, commands, stopCommands)
	for i := 0; i < int(conf.Workers); i++ {
		go handleCommands(conf.Path, commands)
	}
	<-done
	close(stopCommands)
	log.Println("closing consumers...")
	for i := range consumers {
		consumers[i].Close()
	}
}

func handleCommands(path string, commands chan kafka.Command) {
	for c := range commands {
		if err := routeCommand(c, path); err != nil {
			log.Printf("error on routing command: %s", err)
		}
	}
}

func waitForCommands(listen *net.TCPListener, commands chan kafka.Command, stopCommands chan bool) {
	defer close(commands)
	for {
		listen.SetDeadline(time.Now().Add(200 * time.Millisecond))
		conn, err := listen.AcceptTCP()
		if err != nil {
			if softError(err) {
				continue
			}

			log.Printf("unable to accept tcp connection: %s\n", err)
			continue
		}

		if err := conn.SetKeepAlive(true); err != nil {
			log.Printf("unable to set keep alive: %s\n", err)
		}

		go func(conn net.Conn) {
			reader := bufio.NewReader(conn)
			for {
				line, _, err := reader.ReadLine()

				command := kafka.Command{Connection: conn}
				if err == io.EOF {
					command.Type = kafka.TypeClose
					commands <- command
					return
				}

				if err != nil {
					if closedConnection(err) {
						return
					}

					log.Printf("unable to read connection: %s\n", err)
					continue
				}

				if err = json.Unmarshal(line, &command); err != nil {
					log.Printf("error on json: %s\n", err)
					continue
				}
				command.Connection = conn
				select {
				case <-stopCommands:
					return
				case commands <- command:
				}
			}
		}(conn)
	}
}

func routeCommand(c kafka.Command, path string) error {
	commandNames := map[int]string{
		kafka.TypeClose:   "close",
		kafka.TypeConsume: "consume",
		kafka.TypePublish: "publish",
	}

	log.Printf("received command type=%s \n", commandNames[c.Type])
	switch c.Type {
	case kafka.TypePublish:
		var message kafka.Message
		if err := json.Unmarshal([]byte(c.Body), &message); err != nil {
			return err
		}
		return kafka.Publish(c.Connection, message, c.Topic, path)
	case kafka.TypeConsume:
		consumer, err := kafka.NewConsumer(c.ConsumerName, c.Connection, c.Topic, path)
		if err != nil {
			return err
		}
		consumers[consumer.FileName()] = consumer
		go consumer.Start()
		return nil
	case kafka.TypeClose:
		for i := range consumers {
			if consumers[i].Conn == c.Connection {
				consumers[i].Close()
				delete(consumers, i)
			}
		}
		return nil
	}

	return fmt.Errorf("no expected command type: %d\n", c.Type)
}

func softError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
		return true
	}

	// @TODO improve error checking
	return closedConnection(err)
}

func closedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
