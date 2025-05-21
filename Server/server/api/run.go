package api

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	go func() {
		caCert, err := os.ReadFile("../certs/ca.pem")
		if err != nil {
			panic(err)
		}

		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			ClientCAs:  caPool,
			ClientAuth: tls.NoClientCert, // TODO: change back to tls.RequireAndVerifyClientCert
		}

		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.Use(gin.Logger())

		initGetRequests(router)
		initPostRequests(router)

		srv := &http.Server{
			Addr:      ":8080",
			Handler:   router,
			TLSConfig: tlsConfig,
		}

		if err := srv.ListenAndServeTLS("../certs/server.pem", "../certs/server.key"); err != nil {
			return
		}
	}()
}

func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing JWT"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr != string(jwtSecret) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid JWT"})
			return
		}
		c.Next()
	}
}
