package gprc

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/auth"
	"github.com/romanyx/scraper_auth/internal/change"
	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/reset"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/romanyx/scraper_auth/internal/validation"
	"github.com/romanyx/scraper_auth/internal/verify"
	"github.com/romanyx/scraper_auth/proto"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	internalErrMsg = "internal error"
)

// Registrater is a registration service.
type Registrater interface {
	Registrate(context.Context, *reg.Form, *user.User) error
}

// Authenticater logins user.
type Authenticater interface {
	Authenticate(ctx context.Context, email, password string, t *auth.Token) error
}

// Verifier allows to verify user by token.
type Verifier interface {
	Verify(ctx context.Context, token string, u *user.User) error
}

// PwdReseter allows to send user password reset
// instructions.
type PwdReseter interface {
	Reset(ctx context.Context, email string) error
}

// PwdChanger allows to change password.
type PwdChanger interface {
	Change(ctx context.Context, token string, form *change.Form, u *user.User) error
}

// Server is a auth grpc server implementation.
type Server struct {
	RegSrv  Registrater
	AuthSrv Authenticater
	VrfSrv  Verifier
	RstSrv  PwdReseter
	ChgSrv  PwdChanger
}

// NewServer factory build grpc server implementation.
func NewServer(r Registrater, a Authenticater, v Verifier, rst PwdReseter, c PwdChanger) *Server {
	s := Server{
		RegSrv:  r,
		AuthSrv: a,
		VrfSrv:  v,
		RstSrv:  rst,
		ChgSrv:  c,
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
		case validation.Errors:
			return nil, statusFromErrors(v)
		default:
			return nil, status.Error(codes.Internal, internalErrMsg)
		}
	}

	var r proto.UserResponse
	setResp(&r, &u)
	return &r, nil
}

func statusFromErrors(validations validation.Errors) error {
	violations := make([]*epb.QuotaFailure_Violation, len(validations))

	for i, v := range validations {
		violation := epb.QuotaFailure_Violation{
			Subject:     v.Field,
			Description: v.Message,
		}

		violations[i] = &violation
	}

	st := status.New(codes.InvalidArgument, validations.Error())
	ds, err := st.WithDetails(&epb.QuotaFailure{
		Violations: violations,
	})

	if err != nil {
		return st.Err()
	}

	return ds.Err()
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
func (s *Server) Verify(ctx context.Context, req *proto.EmailVerifyRequest) (*proto.UserResponse, error) {
	var u user.User

	if err := s.VrfSrv.Verify(ctx, req.Token, &u); err != nil {
		switch err := errors.Cause(err); err {
		case verify.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, internalErrMsg)
		}
	}

	var r proto.UserResponse
	setResp(&r, &u)
	return &r, nil
}

// Reset send password reset instructions..
func (s *Server) Reset(ctx context.Context, req *proto.PasswordResetRequest) (*proto.PasswordResetResponse, error) {
	if err := s.RstSrv.Reset(ctx, req.Email); err != nil {
		switch err := errors.Cause(err); err {
		case reset.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, internalErrMsg)
		}
	}

	var r proto.PasswordResetResponse
	return &r, nil
}

// Change changes users password.
func (s *Server) Change(ctx context.Context, req *proto.PasswordChangeRequest) (*proto.UserResponse, error) {
	var u user.User
	form := change.Form{
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}

	if err := s.ChgSrv.Change(ctx, req.Token, &form, &u); err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return nil, status.Error(codes.InvalidArgument, v.Error())
		default:
			switch err := errors.Cause(err); err {
			case change.ErrNotFound:
				return nil, status.Error(codes.NotFound, err.Error())
			case change.ErrTokenExpired:
				return nil, status.Error(codes.InvalidArgument, err.Error())
			default:
				return nil, status.Error(codes.Internal, internalErrMsg)
			}
		}
	}

	var r proto.UserResponse
	setResp(&r, &u)
	return &r, nil
}
