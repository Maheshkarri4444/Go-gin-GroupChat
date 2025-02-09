package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	route "githun.com/Maheshkarri4444/group-chat/routes"
)

func main() {
	router := gin.Default()

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true // Be careful with this in production
				},
			},
		},
	})

	// Configure CORS - More permissive for development
	// Update your CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Additional headers for WebSocket
	router.Use(func(c *gin.Context) {
		// Allow WebSocket Upgrade
		if c.Request.Header.Get("Upgrade") == "websocket" {
			c.Next()
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Connection, Upgrade")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// if c.Request.Method == "OPTIONS" {
		// 	c.AbortWithStatus(204)
		// 	return
		// }

		c.Next()
	})

	// Socket.io event handlers
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("New user connected:", s.ID())
		s.Join("chatRoom")
		return nil
	})

	server.OnEvent("/", "sendMessage", func(s socketio.Conn, msg map[string]string) {
		fmt.Println("Message Received:", msg)
		server.BroadcastToRoom("/", "chatRoom", "receiveMessage", msg)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Socket.io error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("User Disconnected:", s.ID(), "Reason:", reason)
	})

	// Start socket.io server
	go func() {
		if err := server.Serve(); err != nil {
			fmt.Println("Socket.io server error:", err)
		}
	}()
	defer server.Close()

	// Socket.io routes
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	// API routes
	authRoutes := router.Group("/api/auth")
	{
		route.AuthRoutes(authRoutes)
	}

	messageRoutes := router.Group("/api/message")
	{
		route.MessageRoutes(messageRoutes)
	}

	// Start the server
	if err := router.Run(":4000"); err != nil {
		fmt.Println("Server error:", err)
	}
}
