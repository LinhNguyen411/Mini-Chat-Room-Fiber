package chat

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	OriginId        string `json:"origin_id"`
	DestinationId   string `json:"destination_id"`
	OriginName      string `json:"origin_name"`
	DestinationName string `json:"destination_name"`
	Content         string `json:"content"`
	Broadcast       bool   `json:"broadcast"`
}

type Client struct {
	Id                 string
	Name               string
	WebsocketConn      *websocket.Conn
	Manager            *ChatManager
	ReceiveMessageChan chan *Message
}

func NewClient(id string, name string, conn *websocket.Conn, manager *ChatManager) *Client {
	return &Client{
		Id:                 id,
		Name:               name,
		WebsocketConn:      conn,
		Manager:            manager,
		ReceiveMessageChan: make(chan *Message),
	}
}

var Wg sync.WaitGroup

func (c *Client) WriteMessages() {
	defer func() {
		Wg.Done()
		c.Manager.UnsubscribeClientChan <- c
		_ = c.WebsocketConn.Close()

		unregisterNotification := &Message{
			OriginId:   "Manager",
			OriginName: "Manager",
			Content:    fmt.Sprintf("*** %s (%s) left this room ***", c.Name, c.Id),
			Broadcast:  true,
		}

		c.Manager.BroadcastNotificationChan <- unregisterNotification
	}()

	for {
		_, msg, err := c.WebsocketConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		chatMessage := Message{}
		json.Unmarshal(msg, &chatMessage)
		chatMessage.OriginId = c.Id
		c.Manager.SendMessageChan <- &chatMessage
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		Wg.Done()
		_ = c.WebsocketConn.Close()
	}()

	for {
		select {
		case messageReceived := <-c.ReceiveMessageChan:
			data, _ := json.Marshal(messageReceived)
			c.WebsocketConn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
