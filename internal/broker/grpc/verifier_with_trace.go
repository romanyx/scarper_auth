package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/trace template

//go:generate gowrap gen -d . -i Verifier -t ../../../templates/trace -o verifier_with_trace.go

import (
	"context"

	"github.com/romanyx/scraper_auth/internal/user"
	"go.opencensus.io/trace"
)

// VerifierWithTracing implements Verifier interface instrumented with opentracing spans
type VerifierWithTracing struct {
	Verifier
}

// NewVerifierWithTracing returns VerifierWithTracing
func NewVerifierWithTracing(base Verifier) VerifierWithTracing {
	d := VerifierWithTracing{
		Verifier: base,
	}

	return d
}

// Verify implements Verifier
func (d VerifierWithTracing) Verify(ctx context.Context, token string, u *user.User) (err error) {
	ctx, span := trace.StartSpan(ctx, "Verifier.Verify")
	defer span.End()

	return d.Verifier.Verify(ctx, token, u)
}
