package authservice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devopshaven/gateway-auth-service/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn       *grpc.ClientConn
	authClient pb.GatewayAuthServiceClient
}

func NewClient(endpoint string) *Client {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))

	if err != nil {
		panic(err)
	}

	client := pb.NewGatewayAuthServiceClient(conn)

	return &Client{
		conn:       conn,
		authClient: client,
	}
}

func (c *Client) Authorize(ctx context.Context, method, host, url string, headers http.Header) (*AuthResult, error) {
	pbh := make([]*pb.Header, 0)

	for k, v := range headers {
		pbh = append(pbh, &pb.Header{
			Name:   k,
			Values: v,
		})
	}

	res, err := c.authClient.Authorize(ctx, &pb.AuthorizeRequest{
		Host:           host,
		Method:         method,
		RequestUrl:     url,
		RequestHeaders: pbh,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot call auth server: %w", err)
	}

	// fmt.Printf("Response from auth server: %+v\n", *res)

	ah := make(http.Header)
	for _, v := range res.AdditionalHeaders {
		for _, val := range v.Values {
			ah.Add(v.Name, val)
		}
	}

	return &AuthResult{
		Block:   res.Block,
		status:  int(res.StatusCode),
		err:     res.ErrorMessage,
		content: res.Content,
		header:  ah,
	}, nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("cannot close auth service client: %w", err)
	}

	return nil
}
