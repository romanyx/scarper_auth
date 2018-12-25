package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/logrus template

//go:generate gowrap gen -d . -i Authenticater -t ../../../templates/logrus -o authenticater_with_logrus.go

import (
	"context"

	"github.com/romanyx/scraper_auth/internal/auth"
	"github.com/sirupsen/logrus"
)

// AuthenticaterWithLogrus implements Authenticater that is instrumented with logrus logger
type AuthenticaterWithLogrus struct {
	log  *logrus.Entry
	base Authenticater
}

// NewAuthenticaterWithLogrus instruments an implementation of the Authenticater with simple logging
func NewAuthenticaterWithLogrus(base Authenticater, log *logrus.Entry) AuthenticaterWithLogrus {
	return AuthenticaterWithLogrus{
		base: base,
		log:  log,
	}
}

// Authenticate implements Authenticater
func (d AuthenticaterWithLogrus) Authenticate(ctx context.Context, email string, password string, t *auth.Token) (err error) {
	d.log.WithFields(logrus.Fields(map[string]interface{}{
		"ctx":      ctx,
		"email":    email,
		"password": password,
		"t":        t})).Debug("AuthenticaterWithLogrus: calling Authenticate")
	defer func() {
		if err != nil {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("AuthenticaterWithLogrus: method Authenticate returned an error")
		} else {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("AuthenticaterWithLogrus: method Authenticate finished")
		}
	}()
	return d.base.Authenticate(ctx, email, password, t)
}
