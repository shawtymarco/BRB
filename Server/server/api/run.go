package api

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"strings"

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

		// Gin setup for API
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.Use(gin.Logger())

		// Apply JWT middleware to your API routes
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
