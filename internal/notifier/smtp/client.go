package smtp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/romanyx/scraper_auth/internal/user"
)

//go:generate go-bindata -prefix templates/ -pkg smtp -o templates.bindata.go templates/

var (
	vrfTmpl = template.Must(template.New("vrf").Parse(string(MustAsset("verify.txt"))))
	cngTmpl = template.Must(template.New("cng").Parse(string(MustAsset("change.txt"))))
)

// Client holds data required to send notifications through email.
type Client struct {
	auth       smtp.Auth
	addr, from string
}

// NewClient factory to initialize smtp client.
func NewClient(auth smtp.Auth, addr, from string) *Client {
	c := Client{
		auth: auth,
		addr: addr,
		from: from,
	}

	return &c
}

// Verify sends confirmation to given emails.
func (c *Client) Verify(ctx context.Context, u *user.User) error {
	to := []string{u.Email}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: Scraper Team <%s>\n", c.from)
	fmt.Fprintf(&buf, "To: %s\n", strings.Join(to, ", "))
	fmt.Fprintf(&buf, "Subject: Email Address Verification\n\n")

	if err := vrfTmpl.Execute(&buf, u); err != nil {
		return errors.Wrap(err, "execute template")
	}

	cn, err := smtp.Dial(c.addr)
	if err != nil {
		return errors.Wrap(err, "dial")
	}

	host, _, _ := net.SplitHostPort(c.addr)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	cn.StartTLS(tlsconfig)

	if err = cn.Auth(c.auth); err != nil {
		return errors.Wrap(err, "auth")
	}

	if err = cn.Mail(c.from); err != nil {
		return errors.Wrap(err, "mail")
	}

	for _, addr := range to {
		if err = cn.Rcpt(addr); err != nil {
			return errors.Wrap(err, "rcpt")
		}
	}

	w, err := cn.Data()
	if err != nil {
		return errors.Wrap(err, "data")
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		return errors.Wrap(err, "write")
	}

	if err := cn.Quit(); err != nil {
		if v, ok := err.(*textproto.Error); ok && v.Code == 250 {
			return nil
		}
		return errors.Wrap(err, "quit")
	}

	return nil
}

// Change sends password change instructions to given emails.
func (c *Client) Change(ctx context.Context, u *user.User, token string) error {
	to := []string{u.Email}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: Scraper Team <%s>\n", c.from)
	fmt.Fprintf(&buf, "To: %s\n", strings.Join(to, ", "))
	fmt.Fprintf(&buf, "Subject: Password Change Instructions\n\n")

	u.Token = &token
	if err := cngTmpl.Execute(&buf, u); err != nil {
		return errors.Wrap(err, "execute template")
	}

	cn, err := smtp.Dial(c.addr)
	if err != nil {
		return errors.Wrap(err, "dial")
	}

	host, _, _ := net.SplitHostPort(c.addr)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	cn.StartTLS(tlsconfig)

	if err = cn.Auth(c.auth); err != nil {
		return errors.Wrap(err, "auth")
	}

	if err = cn.Mail(c.from); err != nil {
		return errors.Wrap(err, "mail")
	}

	for _, addr := range to {
		if err = cn.Rcpt(addr); err != nil {
			return errors.Wrap(err, "rcpt")
		}
	}

	w, err := cn.Data()
	if err != nil {
		return errors.Wrap(err, "data")
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		return errors.Wrap(err, "write")
	}

	if err := cn.Quit(); err != nil {
		if v, ok := err.(*textproto.Error); ok && v.Code == 250 {
			return nil
		}
		return errors.Wrap(err, "quit")
	}

	return nil
}
