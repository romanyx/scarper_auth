package main

import (
	"context"
	"testing"

	"github.com/romanyx/scraper_auth/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestSignUp(t *testing.T) {
	t.Log("with prepared server and client.")
	{
		db, teardown := postgresDB(t)
		defer teardown()
		addr, stop := newServer(db)
		defer stop()
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		cli := proto.NewAuthClient(conn)

		t.Log("\ttest:0\tshould register user.")
		{
			ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
			defer cancel()
			_, err := cli.SignUp(ctx, &proto.SignUpRequest{
				Email:                "john@example.com",
				AccountId:            "492c9a6d-255e-4a61-a460-2d622d4c6e96",
				Password:             "password",
				PasswordConfirmation: "password",
			})
			assert.Nil(t, err)
		}
	}
}
