package api

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

		gin.SetMode(gin.ReleaseMode)
		router := gin.New()
		router.Use(gin.Recovery())

		// custom logger that skips /api/connect
		router.Use(func(c *gin.Context) {
			if c.Request.URL.Path == "/api/connect" {
				c.Next()
				return
			}

			start := time.Now()
			c.Next()
			latency := time.Since(start)
			status := c.Writer.Status()
			log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s  %s\n",
				start.Format("2006/01/02 - 15:04:05"),
				status,
				latency,
				c.ClientIP(),
				c.Request.Method,
				c.Request.URL.Path,
			)
		})

		// Apply JWT middleware
		apiGroup := router.Group("/api")
		apiGroup.Use(jwtAuthMiddleware())
		initGetRequests(apiGroup)
		initPostRequests(apiGroup)

		srv := &http.Server{
			Addr:      ":8080",
			Handler:   router,
			TLSConfig: tlsConfig,
		}

		log.Println(`Starting secure server with API on :8080`)
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
