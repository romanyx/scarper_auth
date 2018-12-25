package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/logrus template

//go:generate gowrap gen -d . -i PwdChanger -t ../../../templates/logrus -o pwd_changer_with_logrus.go

import (
	"context"

	"github.com/romanyx/scraper_auth/internal/change"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/sirupsen/logrus"
)

// PwdChangerWithLogrus implements PwdChanger that is instrumented with logrus logger
type PwdChangerWithLogrus struct {
	log  *logrus.Entry
	base PwdChanger
}

// NewPwdChangerWithLogrus instruments an implementation of the PwdChanger with simple logging
func NewPwdChangerWithLogrus(base PwdChanger, log *logrus.Entry) PwdChangerWithLogrus {
	return PwdChangerWithLogrus{
		base: base,
		log:  log,
	}
}

// Change implements PwdChanger
func (d PwdChangerWithLogrus) Change(ctx context.Context, token string, form *change.Form, u *user.User) (err error) {
	d.log.WithFields(logrus.Fields(map[string]interface{}{
		"ctx":   ctx,
		"token": token,
		"form":  form,
		"u":     u})).Debug("PwdChangerWithLogrus: calling Change")
	defer func() {
		if err != nil {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("PwdChangerWithLogrus: method Change returned an error")
		} else {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("PwdChangerWithLogrus: method Change finished")
		}
	}()
	return d.base.Change(ctx, token, form, u)
}
