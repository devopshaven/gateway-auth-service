package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authservice "github.com/devopshaven/gateway-auth-service"
	"github.com/google/uuid"
)

// random generator
var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

// Auth request handler function
func handleAuthRequest(req *authservice.AuthRequest) error {
	// fmt.Printf("auth request received! yipeee 🤓: %+v\n", *req)

	if r.Intn(2) != 0 {
		// Block request example
		req.BlockWithError(500, "Request blocked")

		return nil
	}

	// Allow request with additional headers.
	headers := make(http.Header)

	// These two headers will also passed to the upstream server.
	headers.Set("X-Auth-User-Id", uuid.NewString())
	headers.Set("X-Laos", "panda")

	req.Passtrough(headers)

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Creating new instance from the auth service
	authservice.NewGatewayAuthService(ctx, "localhost:5009", handleAuthRequest)

	fmt.Println("Auth service stared on port: 5009 (gRPC)")

	// Wait for CTRL+C or terminate/interrupt os signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("Server terminated")
}
