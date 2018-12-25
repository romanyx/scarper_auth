package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/trace template

//go:generate gowrap gen -d . -i PwdReseter -t ../../../templates/trace -o pwd_reseter_with_trace.go

import (
	"context"

	"go.opencensus.io/trace"
)

// PwdReseterWithTracing implements PwdReseter interface instrumented with opentracing spans
type PwdReseterWithTracing struct {
	PwdReseter
}

// NewPwdReseterWithTracing returns PwdReseterWithTracing
func NewPwdReseterWithTracing(base PwdReseter) PwdReseterWithTracing {
	d := PwdReseterWithTracing{
		PwdReseter: base,
	}

	return d
}

// Reset implements PwdReseter
func (d PwdReseterWithTracing) Reset(ctx context.Context, email string) (err error) {
	ctx, span := trace.StartSpan(ctx, "PwdReseter.Reset")
	defer span.End()

	return d.PwdReseter.Reset(ctx, email)
}
