# gateway-auth-service

[![Go Reference](https://pkg.go.dev/badge/github.com/devopshaven/gateway-auth-service.svg)](https://pkg.go.dev/github.com/devopshaven/gateway-auth-service)

Authorization service for [DevopsHaven API Gateway](https://github.com/devopshaven/api-gateway)

This snippet shows how can you implement your authorization logic for your services.

The service communicates with the gateway via insecure gRPC protocol.

```go
package main

import authservice "github.com/devopshaven/gateway-auth-service"

...

// Auth request handler function
func handleAuthRequest(req *authservice.AuthRequest) error {
	// Allow request with additional headers.
	headers := make(http.Header)

	// These two headers will also passed to the upstream server.
	headers.Set("X-Auth-User-Id", uuid.NewString())
	headers.Set("X-Laos", "panda")

    // Allow request
	req.Passtrough(headers)

    // or Block request example with message
	req.BlockWithError(500, "Request blocked")

	return nil
}

func main() {
    ...

	// Creating new instance from the auth service and listen on port 5009
	authservice.NewGatewayAuthService(ctx, ":5009", handleAuthRequest)

	...
}
```