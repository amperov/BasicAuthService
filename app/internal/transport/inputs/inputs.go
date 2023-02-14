package inputs

import "GGSAuth/app/internal/transport/grpc"

type SignRequest struct {
	Email    string
	Password string
}

func (r *SignRequest) UpFromGRPC(request *grpc.SignUpRequest) {
	r.Email = request.Email
	r.Password = request.Password
}

func (r *SignRequest) InFromGRPC(request *grpc.SignInRequest) {
	r.Email = request.Email
	r.Password = request.Password
}
