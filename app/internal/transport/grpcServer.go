package transport

import (
	"context"
	"github.com/amperov/basic-auth-service/app/internal/transport/grpc"
)

type AuthService interface {
	SignUp(ctx context.Context, email, password string) (int, string, string, error) //UserID, AccessCode, Status, Error
	SignIn(ctx context.Context, email, password string) (int, string, string, error)
	Identify(ctx context.Context, AccessCode string) (int, string, error)
}
type GRPCServer struct {
	AuthService AuthService
	grpc.UnimplementedAuthorizationServer
}

func NewGRPCServer(authService AuthService) *GRPCServer {
	return &GRPCServer{AuthService: authService}
}

func (s *GRPCServer) mustEmbedUnimplementedAuthorizationServer() {
	//TODO implement me
	panic("implement me")
}

func (s *GRPCServer) SignUp(ctx context.Context, request *grpc.SignUpRequest) (*grpc.SignResponse, error) {
	var Request inputs.SignRequest
	Request.UpFromGRPC(request)

	UserID, AccessCode, Status, err := s.AuthService.SignUp(ctx, Request.Email, Request.Password)
	if err != nil {
		return nil, err
	}

	return &grpc.SignResponse{UserID: int64(UserID), AccessCode: AccessCode, Status: Status}, nil
}

func (s *GRPCServer) SignIn(ctx context.Context, request *grpc.SignInRequest) (*grpc.SignResponse, error) {
	var Request inputs.SignRequest
	Request.InFromGRPC(request)

	return &grpc.SignResponse{}, nil
}

func (s *GRPCServer) Identity(ctx context.Context, request *grpc.IdentityRequest) (*grpc.IdentityResponse, error) {
	UserID, AccessCode, err := s.AuthService.Identify(ctx, request.GetAccessCode())
	if err != nil {
		return nil, err
	}

	return &grpc.IdentityResponse{
		UserID:     int64(UserID),
		Status:     err.Error(),
		AccessCode: AccessCode,
	}, nil
}

func (s *GRPCServer) RecoverPassword(ctx context.Context, request *grpc.RecoverPasswordRequest) (*grpc.RecoverPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GRPCServer) ChangePassword(ctx context.Context, request *grpc.ChangePasswordRequest) (*grpc.ChangePasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GRPCServer) AcceptAction(ctx context.Context, request *grpc.AcceptActionRequest) (*grpc.AcceptActionResponse, error) {
	//TODO implement me
	panic("implement me")
}
