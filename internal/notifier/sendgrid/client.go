package sendgrid

import (
	"bytes"
	"context"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/romanyx/scraper_auth/internal/user"
)

//go:generate go-bindata -prefix templates/ -pkg sendgrid -o templates.bindata.go templates/

const (
	verifyTitle = "Email Address Verification"
	changeTitle = "Subject: Password Change Instructions"
)

var (
	vrfTmpl = template.Must(template.New("vrf").Parse(string(MustAsset("verify.txt"))))
	cngTmpl = template.Must(template.New("cng").Parse(string(MustAsset("change.txt"))))
)

// Client holds data required to send notifications through email.
type Client struct {
	client *sendgrid.Client
	from   *mail.Email
}

// NewClient factory to initialize smtp client.
func NewClient(cli *sendgrid.Client, from *mail.Email) *Client {
	c := Client{
		client: cli,
		from:   from,
	}

	return &c
}

// Verify sends confirmation to given emails.
func (c *Client) Verify(ctx context.Context, u *user.User) error {
	to := []string{u.Email}

	var buf bytes.Buffer
	if err := vrfTmpl.Execute(&buf, u); err != nil {
		return errors.Wrap(err, "execute template")
	}

	content := mail.NewContent("text/html", buf.String())
	for _, addr := range to {
		toEmail := mail.NewEmail("", addr)
		message := mail.NewV3MailInit(c.from, verifyTitle, toEmail, content)

		_, err := c.client.Send(message)
		if err != nil {
			return errors.Wrap(err, "verify email send")
		}
	}

	return nil
}

// Change sends password change instructions to given emails.
func (c *Client) Change(ctx context.Context, u *user.User, token string) error {
	to := []string{u.Email}

	var buf bytes.Buffer
	u.Token = &token
	if err := cngTmpl.Execute(&buf, u); err != nil {
		return errors.Wrap(err, "execute template")
	}

	content := mail.NewContent("text/html", buf.String())
	for _, addr := range to {
		toEmail := mail.NewEmail("", addr)
		message := mail.NewV3MailInit(c.from, changeTitle, toEmail, content)

		_, err := c.client.Send(message)
		if err != nil {
			return errors.Wrap(err, "change email send")
		}

	}

	return nil
}
