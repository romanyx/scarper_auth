package gprc

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/auth"
	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/romanyx/scraper_auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	internalErrMsg = "internal error"
)

// Registrater is a registration serviceration service.
type Registrater interface {
	Registrate(context.Context, *reg.Form, *user.User) error
}

// Authenticater logins user.
type Authenticater interface {
	Authenticate(ctx context.Context, email, password string, t *auth.Token) error
}

// Server is a auth grpc server implementation.
type Server struct {
	RegSrv  Registrater
	AuthSrv Authenticater
}

// NewServer factory build grpc server implementation.
func NewServer(r Registrater, a Authenticater) *Server {
	s := Server{
		RegSrv:  r,
		AuthSrv: a,
	}

	return &s
}

// SignUp registration implementation.
func (s *Server) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.UserResponse, error) {
	var f reg.Form
	var u user.User
	setForm(req, &f)
	if err := s.RegSrv.Registrate(ctx, &f, &u); err != nil {
		switch v := errors.Cause(err).(type) {
		case reg.ValidationErrors:
			return nil, status.Error(codes.InvalidArgument, v.Error())
		default:
			fmt.Println(err)
			return nil, status.Error(codes.Internal, internalErrMsg)
		}
	}

	var r proto.UserResponse
	setResp(&r, &u)
	return &r, nil
}

func setForm(req *proto.SignUpRequest, f *reg.Form) {
	f.Email = req.Email
	f.AccountID = req.AccountId
	f.Password = req.Password
	f.PasswordConfirmation = req.PasswordConfirmation
}

func setResp(resp *proto.UserResponse, u *user.User) {
	resp.AccountId = u.AccountID
	resp.Status = proto.UserStatus(proto.UserStatus_value[u.Status])
	resp.Email = u.Email
	ca, _ := ptypes.TimestampProto(u.CreatedAt)
	resp.CreatedAt = ca
	ua, _ := ptypes.TimestampProto(u.UpdatedAt)
	resp.UpdatedAt = ua
}

// SignIn authentication implementation.
func (s *Server) SignIn(ctx context.Context, req *proto.SignInRequest) (*proto.SignInResponse, error) {
	var t auth.Token

	if err := s.AuthSrv.Authenticate(ctx, req.Email, req.Password, &t); err != nil {
		switch err := errors.Cause(err); err {
		case auth.ErrNotFound, auth.ErrWrongPassword, auth.ErrNotVerified:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, internalErrMsg)
		}
	}

	resp := proto.SignInResponse{
		Token: t.Token,
	}

	return &resp, nil
}

// Verify verifies users email.
func (s *Server) Verify(context.Context, *proto.EmailVerifyRequest) (*proto.UserResponse, error) {
	return nil, nil
}

// Reset resets users password.
func (s *Server) Reset(context.Context, *proto.PasswordResetRequest) (*proto.PasswordResetResponse, error) {
	return nil, nil
}
