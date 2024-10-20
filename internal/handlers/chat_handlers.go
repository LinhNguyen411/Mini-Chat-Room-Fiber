package handlers

import (
	"fmt"

	"github.com/LinhNguyen411/chat-room-fiber/internal/chat"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterRoomViewHandler(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodPost {

		nickName := c.FormValue("nick")
		return c.Redirect(fmt.Sprintf("/room/%s", nickName))
	}

	return c.Render("internal/views/register.html", nil)
}

func ChatRoomViewHandler(c *fiber.Ctx) error {
	data := fiber.Map{
		"nick": c.Params("nick"),
	}
	return c.Render("internal/views/room.html", data)
}

func RegisterHandler(c *websocket.Conn) {
	chat.Wg.Add(2)

	client := chat.NewClient(uuid.New().String(), c.Params("nick"), c, &chat.Manager)
	client.Manager.SubscribeClientChan <- client

	registerNotification := &chat.Message{
		OriginId:   "Manager",
		OriginName: "Manager",
		Content:    fmt.Sprintf("***  %s (%s) joined to this room ***", client.Name, client.Id),
		Broadcast:  true,
	}
	client.Manager.BroadcastNotificationChan <- registerNotification

	go client.ReadMessages()
	go client.WriteMessages()

	chat.Wg.Wait()
}
