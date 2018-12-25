package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/logrus template

//go:generate gowrap gen -d . -i Registrater -t ../../../templates/logrus -o registrater_with_logrus.go

import (
	"context"

	"github.com/romanyx/scraper_auth/internal/reg"
	"github.com/romanyx/scraper_auth/internal/user"
	"github.com/sirupsen/logrus"
)

// RegistraterWithLogrus implements Registrater that is instrumented with logrus logger
type RegistraterWithLogrus struct {
	log  *logrus.Entry
	base Registrater
}

// NewRegistraterWithLogrus instruments an implementation of the Registrater with simple logging
func NewRegistraterWithLogrus(base Registrater, log *logrus.Entry) RegistraterWithLogrus {
	return RegistraterWithLogrus{
		base: base,
		log:  log,
	}
}

// Registrate implements Registrater
func (d RegistraterWithLogrus) Registrate(ctx context.Context, fp1 *reg.Form, up1 *user.User) (err error) {
	d.log.WithFields(logrus.Fields(map[string]interface{}{
		"ctx": ctx,
		"fp1": fp1,
		"up1": up1})).Debug("RegistraterWithLogrus: calling Registrate")
	defer func() {
		if err != nil {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("RegistraterWithLogrus: method Registrate returned an error")
		} else {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("RegistraterWithLogrus: method Registrate finished")
		}
	}()
	return d.base.Registrate(ctx, fp1, up1)
}
