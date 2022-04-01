package authservice

import (
	"context"
	"net"
	"net/http"

	"github.com/devopshaven/gateway-auth-service/internal/pb"
	"google.golang.org/grpc"
)

// AuthRequest authorization request from the Gateway service.
// You can block or allow pass-trough with Passtrough and BlockWithContent or BlockWithError methods.
type AuthRequest struct {
	ContentType string
	Method      string
	Host        string
	URL         string
	Header      http.Header

	res pb.AuthorizeResponse
}

// Passtrough allow the gateway to send the data to the upstream.
func (ar *AuthRequest) Passtrough(additionalHeaders http.Header) {
	var headers []*pb.Header

	// Transform headers slice to http.Header
	for k, v := range additionalHeaders {
		headers = append(headers, &pb.Header{
			Name:   k,
			Values: v,
		})
	}

	ar.res.Block = false
	ar.res.AdditionalHeaders = headers
}

// BlockWithError blocks the request with specific content.
func (ar *AuthRequest) BlockWithContent(status int, contentType string, content []byte) {
	ar.res.Block = true
	ar.res.ReplaceContent = true
	ar.res.StatusCode = int32(status)
	ar.res.Content = content

	ar.res.AdditionalHeaders = append(ar.res.AdditionalHeaders, &pb.Header{
		Name:   "Content-Type",
		Values: []string{contentType},
	})
}

// BlockWithError block the request with plain text error.
func (ar *AuthRequest) BlockWithError(status int, message string) {
	ar.res.Block = true
	ar.res.ReplaceContent = false
	ar.res.StatusCode = int32(status)
	ar.res.ErrorMessage = message
}

func newAuthRequest(req *pb.AuthorizeRequest) *AuthRequest {
	ct, _ := req.ContentType()

	headers := http.Header{}

	// Transform headers slice to http.Header
	for _, h := range req.RequestHeaders {
		headers[h.Name] = h.Values
	}

	return &AuthRequest{
		ContentType: ct,
		Method:      req.Method,
		Host:        req.Host,
		URL:         req.RequestUrl,
		Header:      headers,

		res: pb.AuthorizeResponse{
			StatusCode:        0,
			ErrorMessage:      "",
			RemoveHeaders:     []string{},
			AdditionalHeaders: []*pb.Header{},
			ReplaceContent:    false,
		},
	}
}

type AuthHandler func(*AuthRequest) error

type server struct {
	pb.UnimplementedGatewayAuthServiceServer

	handler AuthHandler
}

// Authorize the implemented gRPC method which will call the AuthHandler callback.
func (s *server) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	ar := newAuthRequest(req)
	s.handler(ar)

	return &ar.res, nil
}

type AuthService struct {
}

// NewGatewayAuthService initializes a new instance from the auth service. The third parameter is the callback function for authentication requests.
// There you can block or allow pass-trough requests with additional header manipulation.
func NewGatewayAuthService(ctx context.Context, listen string, handler AuthHandler) {

	s := grpc.NewServer()

	srv := &server{
		handler: handler,
	}
	pb.RegisterGatewayAuthServiceServer(s, srv)

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		panic(err)
	}

	go func() {
		go s.Serve(lis)

		for range ctx.Done() {
			s.Stop()
		}
	}()
}
