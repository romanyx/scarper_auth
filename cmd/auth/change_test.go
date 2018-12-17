package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/romanyx/polluter"
	"github.com/romanyx/scraper_auth/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	changeData = `
users:
  - id: 1
    email: john@example.com
    account_id: 492c9a6d-255e-4a61-a460-2d622d4c6e96
    password_hash: $2a$10$nhWT4xXRkk0aoqOMOs7UyOBJ1f/XXFGt5rYDBo9CAQnruBvg.U3d6
    status: VERIFIED
    token: 492c9a6d-255e-4a61-a460-2d622d4c6e96
resets:
  - user_id: 1
    token: 492c9a6d-255e-4a61-a460-2d622d4c6e96
    expired_at: %s
`
)

func TestChange(t *testing.T) {
	t.Log("with prepared server and client.")
	{
		db, teardown := postgresDB(t)
		defer teardown()
		p := polluter.New(polluter.PostgresEngine(db))
		tm := time.Now().Add(10 * time.Hour)
		assert.Nil(t, p.Pollute(strings.NewReader(fmt.Sprintf(changeData, tm.Format(time.RFC3339)))))

		addr, stop := newServer(db)
		defer stop()
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		cli := proto.NewAuthClient(conn)

		t.Log("\ttest:0\tshould change user password.")
		{
			ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
			defer cancel()
			resp, err := cli.Change(ctx, &proto.PasswordChangeRequest{
				Token:                "492c9a6d-255e-4a61-a460-2d622d4c6e96",
				Password:             "password",
				PasswordConfirmation: "password",
			})
			assert.Nil(t, err)
			assert.NotNil(t, resp)
		}
	}
}
