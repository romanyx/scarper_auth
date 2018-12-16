package main

import (
	"context"
	"strings"
	"testing"

	"github.com/romanyx/polluter"
	"github.com/romanyx/scraper_auth/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	userData = `
users:
  - email: john@example.com
    account_id: 492c9a6d-255e-4a61-a460-2d622d4c6e96
    password_hash: $2a$10$nhWT4xXRkk0aoqOMOs7UyOBJ1f/XXFGt5rYDBo9CAQnruBvg.U3d6
    status: VERIFIED
    token: 492c9a6d-255e-4a61-a460-2d622d4c6e96
`
)

func TestSignIn(t *testing.T) {
	t.Log("with prepared server and client.")
	{
		db, teardown := postgresDB(t)
		defer teardown()
		p := polluter.New(polluter.PostgresEngine(db))
		assert.Nil(t, p.Pollute(strings.NewReader(userData)))

		addr, stop := newServer(db)
		defer stop()
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		cli := proto.NewAuthClient(conn)

		t.Log("\tTest:0\tshould sign in user.")
		{
			ctx, cancel := context.WithTimeout(context.Background(), caseTiemout)
			defer cancel()
			resp, err := cli.SignIn(ctx, &proto.SignInRequest{
				Email:    "john@example.com",
				Password: "password",
			})
			assert.Nil(t, err)
			assert.NotNil(t, resp)

			claims, err := a.ParseClaims(ctx, resp.Token)
			assert.Nil(t, err)
			assert.Equal(t, claims.Subject, "492c9a6d-255e-4a61-a460-2d622d4c6e96")
		}
	}
}
