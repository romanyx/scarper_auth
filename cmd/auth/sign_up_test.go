package main

import (
	"context"
	"testing"
	"time"

	"github.com/romanyx/scraper_auth/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	addr        = "localhost:50051"
	caseTiemout = 3 * time.Second
)

func TestSignUp(t *testing.T) {
	t.Log("with prepared server and client.")
	{
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		cli := proto.NewAuthClient(conn)

		t.Log("\tTest:0\tshould register user.")
		{
			ctx, cancel := context.WithTimeout(context.Background(), caseTiemout)
			defer cancel()
			_, err := cli.SignUp(ctx, &proto.SignUpRequest{
				Email:                "john@example.com",
				Password:             "password",
				PasswordConfirmation: "password_confirmation",
			})
			assert.Nil(t, err)
		}
	}
}
