package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/trace template

//go:generate gowrap gen -d . -i Authenticater -t ../../../templates/trace -o authenticater_with_trace.go

import (
	"context"

	"github.com/romanyx/scraper_auth/internal/auth"
	"go.opencensus.io/trace"
)

// AuthenticaterWithTracing implements Authenticater interface instrumented with opentracing spans
type AuthenticaterWithTracing struct {
	Authenticater
}

// NewAuthenticaterWithTracing returns AuthenticaterWithTracing
func NewAuthenticaterWithTracing(base Authenticater) AuthenticaterWithTracing {
	d := AuthenticaterWithTracing{
		Authenticater: base,
	}

	return d
}

// Authenticate implements Authenticater
func (d AuthenticaterWithTracing) Authenticate(ctx context.Context, email string, password string, t *auth.Token) (err error) {
	ctx, span := trace.StartSpan(ctx, "Authenticater.Authenticate")
	defer span.End()

	return d.Authenticater.Authenticate(ctx, email, password, t)
}
