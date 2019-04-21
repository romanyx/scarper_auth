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

func TestReset(t *testing.T) {
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

		t.Log("\ttest:0\tcreate reset token.")
		{
			ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
			defer cancel()
			resp, err := cli.Reset(ctx, &proto.PasswordResetRequest{
				Email: "work@romanyx.ru",
			})
			assert.Nil(t, err)
			assert.NotNil(t, resp)
		}
	}
}
