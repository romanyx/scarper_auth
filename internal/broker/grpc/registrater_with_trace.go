package gprc

import (
	"context"

	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/user"
	"go.opencensus.io/trace"
)

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/trace template

//go:generate gowrap gen -d . -i Registrater -t ../../../templates/trace -o registrater_with_trace.go

// RegistraterWithTracing implements Registrater interface instrumented with opentracing spans
type RegistraterWithTracing struct {
	Registrater
}

// NewRegistraterWithTracing returns RegistraterWithTracing
func NewRegistraterWithTracing(base Registrater) RegistraterWithTracing {
	d := RegistraterWithTracing{
		Registrater: base,
	}

	return d
}

// Registrate implements Registrater
func (d RegistraterWithTracing) Registrate(ctx context.Context, fp1 *reg.Form, up1 *user.User) (err error) {
	ctx, span := trace.StartSpan(ctx, "Registrater.Registrate")
	defer span.End()

	return d.Registrater.Registrate(ctx, fp1, up1)
}
