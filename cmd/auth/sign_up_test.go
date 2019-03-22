package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/romanyx/scraper_auth/internal/validation"
	"github.com/romanyx/scraper_auth/proto"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
		t.Log("\ttest:1\tshould return validation errors.")
		{
			ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
			defer cancel()
			_, err := cli.SignUp(ctx, &proto.SignUpRequest{
				Email:                "",
				Password:             "password",
				PasswordConfirmation: "missmatch",
			})

			if err == nil {
				t.Errorf("expected validation error")
			}

			st := status.Convert(err)

			if st.Code() != codes.InvalidArgument {
				t.Errorf("unexpected code: %d expected: %d", st.Code(), codes.InvalidArgument)
			}

			expect := validation.Errors{
				{
					Field:   "email",
					Message: "cannot be blank",
				},
				{
					Field:   "account_id",
					Message: "cannot be blank",
				},
				{
					Field:   "password_confirmation",
					Message: "mismatch",
				},
				{
					Field:   "password",
					Message: "mismatch",
				},
			}

			got := detailsToValidations(st.Details())

			if !reflect.DeepEqual(expect, got) {
				t.Errorf("expected errors to be: %#+v got: %#+v", expect, got)
			}
		}
	}
}

func detailsToValidations(details []interface{}) validation.Errors {
	validations := make(validation.Errors, 0)
	for _, d := range details {
		if f, ok := d.(*epb.QuotaFailure); ok {
			for _, v := range f.Violations {
				validations = append(validations, validation.Error{
					Field:   v.Subject,
					Message: v.Description,
				})
			}
		}
	}

	return validations
}
