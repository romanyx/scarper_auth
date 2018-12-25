package gprc

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../templates/logrus template

//go:generate gowrap gen -d . -i PwdReseter -t ../../../templates/logrus -o pwd_reseter_with_logrus.go

import (
	"context"

	"github.com/sirupsen/logrus"
)

// PwdReseterWithLogrus implements PwdReseter that is instrumented with logrus logger
type PwdReseterWithLogrus struct {
	log  *logrus.Entry
	base PwdReseter
}

// NewPwdReseterWithLogrus instruments an implementation of the PwdReseter with simple logging
func NewPwdReseterWithLogrus(base PwdReseter, log *logrus.Entry) PwdReseterWithLogrus {
	return PwdReseterWithLogrus{
		base: base,
		log:  log,
	}
}

// Reset implements PwdReseter
func (d PwdReseterWithLogrus) Reset(ctx context.Context, email string) (err error) {
	d.log.WithFields(logrus.Fields(map[string]interface{}{
		"ctx":   ctx,
		"email": email})).Debug("PwdReseterWithLogrus: calling Reset")
	defer func() {
		if err != nil {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("PwdReseterWithLogrus: method Reset returned an error")
		} else {
			d.log.WithFields(logrus.Fields(map[string]interface{}{
				"err": err})).Debug("PwdReseterWithLogrus: method Reset finished")
		}
	}()
	return d.base.Reset(ctx, email)
}
