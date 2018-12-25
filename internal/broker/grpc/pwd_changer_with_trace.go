package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/trace template

//go:generate gowrap gen -d . -i PwdChanger -t ../../../templates/trace -o pwd_changer_with_trace.go

import (
	"context"

	"github.com/romanyx/scraper_auth/internal/change"
	"github.com/romanyx/scraper_auth/internal/user"
	"go.opencensus.io/trace"
)

// PwdChangerWithTracing implements PwdChanger interface instrumented with opentracing spans
type PwdChangerWithTracing struct {
	PwdChanger
}

// NewPwdChangerWithTracing returns PwdChangerWithTracing
func NewPwdChangerWithTracing(base PwdChanger) PwdChangerWithTracing {
	d := PwdChangerWithTracing{
		PwdChanger: base,
	}

	return d
}

// Change implements PwdChanger
func (d PwdChangerWithTracing) Change(ctx context.Context, token string, form *change.Form, u *user.User) (err error) {
	ctx, span := trace.StartSpan(ctx, "PwdChanger.Change")
	defer span.End()

	return d.PwdChanger.Change(ctx, token, form, u)
}
