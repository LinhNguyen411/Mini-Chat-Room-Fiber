package main

import (
	"github.com/LinhNguyen411/chat-room-fiber/internal/chat"
	"github.com/LinhNguyen411/chat-room-fiber/internal/handlers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	go chat.Manager.Start()

	app.Get("/api/ws/register/:nick", websocket.New(handlers.RegisterHandler))

	app.Get("/room/:nick", handlers.ChatRoomViewHandler)
	app.All("/", handlers.RegisterRoomViewHandler)

	app.Listen("127.0.0.1:3000")
}
