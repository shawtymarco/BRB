package api

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/server/utils"
	"server/server/ws"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	go func() {
		// Load CA cert
		caCert, err := os.ReadFile("../certs/ca.pem")
		if err != nil {
			panic(err)
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caCert) {
			panic("Failed to append CA cert")
		}

		tlsConfig := &tls.Config{
			ClientCAs: caPool,
			// Change as needed: tls.NoClientCert or tls.RequireAndVerifyClientCert
			ClientAuth: tls.RequireAndVerifyClientCert,
			MinVersion: tls.VersionTLS12,
		}

		// Gin setup for API
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.Use(gin.Logger())

		// Apply JWT middleware to your API routes
		apiGroup := router.Group("/api")
		apiGroup.Use(jwtAuthMiddleware())
		initGetRequests(apiGroup)
		initPostRequests(apiGroup)

		// Add WebSocket route to same router, handle auth manually in WS handler
		router.GET("/ws", wsHandler)

		srv := &http.Server{
			Addr:      ":8080",
			Handler:   router,
			TLSConfig: tlsConfig,
		}

		log.Println("Starting secure server with API + WebSocket on https://localhost:8080")
		if err := srv.ListenAndServeTLS("../certs/server.pem", "../certs/server.key"); err != nil {
			log.Fatal(err)
		}
	}()
}

// JWT middleware for Gin API routes
func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing JWT"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr != string(jwtSecret) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			return
		}
		c.Next()
	}
}

// Gorilla websocket upgrader (allow all origins or customize)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Or restrict origins here
	},
}

func wsHandler(c *gin.Context) {
	// Manually validate JWT from header (since Gin middleware won't be applied here)
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing JWT"})
		return
	}
	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	if tokenStr != string(jwtSecret) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var response ws.DiscordCommand
		utils.Panic(json.Unmarshal(message, &response))

		log.Printf("WS Received: %s", message)

		fmt.Println(response)
		success := "ok"

		utils.Panic(conn.WriteMessage(mt, []byte(success)))
	}
}
