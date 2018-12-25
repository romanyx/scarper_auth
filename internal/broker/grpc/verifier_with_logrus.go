package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/logrus template

//go:generate gowrap gen -d . -i Verifier -t ../../../templates/logrus -o verifier_with_logrus.go

import (
	"context"

	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/sirupsen/logrus"
)

// VerifierWithLogrus implements Verifier that is instrumented with logrus logger
type VerifierWithLogrus struct {
	log  *logrus.Entry
	base Verifier
}

// NewVerifierWithLogrus instruments an implementation of the Verifier with simple logging
func NewVerifierWithLogrus(base Verifier, log *logrus.Entry) VerifierWithLogrus {
	return VerifierWithLogrus{
		base: base,
		log:  log,
	}
}

// Verify implements Verifier
func (d VerifierWithLogrus) Verify(ctx context.Context, token string, u *user.User) (err error) {
	d.log.WithFields(logrus.Fields(map[string]interface{}{
		"ctx":   ctx,
		"token": token,
		"u":     u})).Debug("VerifierWithLogrus: calling Verify")
	defer func() {
		if err != nil {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("VerifierWithLogrus: method Verify returned an error")
		} else {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("VerifierWithLogrus: method Verify finished")
		}
	}()
	return d.base.Verify(ctx, token, u)
}
