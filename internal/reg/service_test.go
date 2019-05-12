package reg

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	user "github.com/romanyx/scraper_auth/internal/user"
	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
	commit := func() error {
		return nil
	}
	rollback := func() error {
		return nil
	}
	tests := []struct {
		name           string
		validateFunc   func(context.Context, *Form) error
		repositoryFunc func(mock *MockRepository)
		informFunc     func(context.Context, *user.User) error
		wantErr        bool
	}{
		{
			name: "ok",
			validateFunc: func(context.Context, *Form) error {
				return nil
			},
			informFunc: func(context.Context, *user.User) error {
				return nil
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(commit, nil, nil)
			},
		},
		{
			name: "validation",
			validateFunc: func(context.Context, *Form) error {
				return errors.New("mock error")
			},
			repositoryFunc: func(m *MockRepository) {},
			wantErr:        true,
		},
		{
			name: "create",
			validateFunc: func(context.Context, *Form) error {
				return nil
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "inform",
			validateFunc: func(context.Context, *Form) error {
				return nil
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, rollback, nil)
			},
			informFunc: func(context.Context, *user.User) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := NewMockRepository(ctrl)
			tt.repositoryFunc(repo)

			s := NewService(repo, validaterFunc(tt.validateFunc), informerFunc(tt.informFunc))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var usr user.User
			form := Form{
				Email:     "john@example.com",
				AccountID: "account_id",
			}
			err := s.Registrate(ctx, &form, &usr)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

type informerFunc func(context.Context, *user.User) error

func (f informerFunc) Verify(ctx context.Context, u *user.User) error {
	return f(ctx, u)
}

type validaterFunc func(context.Context, *Form) error

func (vf validaterFunc) Validate(ctx context.Context, f *Form) error {
	return vf(ctx, f)
}
