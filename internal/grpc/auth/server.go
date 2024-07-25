package auth

import (
	"context"

	ssov1 "github.com/sollidy/go-sso-protos/gen/go/proto/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	// ssov1.UnimplementedAuthServiceServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{
		Token: "token sample",
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("implement")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("implement")
}
