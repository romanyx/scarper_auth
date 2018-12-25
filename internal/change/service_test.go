package change

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	user "github.com/romanyx/scraper_auth/internal/user"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Change(t *testing.T) {
	tests := []struct {
		name         string
		repoFunc     func(m *MockRepository)
		validateFunc func(context.Context, *Form) error
		wantErr      bool
	}{
		{
			name: "ok",
			repoFunc: func(m *MockRepository) {
				call := m.EXPECT().FindResetToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				call.Do(func(ctx context.Context, token string, t *Token) error {
					t.ExpiredAt = time.Now().Add(2 * time.Hour)
					return nil
				})
				call.Return(nil)
				m.EXPECT().ChangePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			validateFunc: func(context.Context, *Form) error {
				return nil
			},
		},
		{
			name: "token",
			repoFunc: func(m *MockRepository) {
				call := m.EXPECT().FindResetToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				call.Do(func(ctx context.Context, token string, t *Token) error {
					t.ExpiredAt = time.Now()
					return nil
				})
				call.Return(nil)
			},
			validateFunc: func(context.Context, *Form) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "validater",
			repoFunc: func(m *MockRepository) {
				call := m.EXPECT().FindResetToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				call.Do(func(ctx context.Context, token string, t *Token) error {
					t.ExpiredAt = time.Now().Add(2 * time.Hour)
					return nil
				})
				call.Return(nil)
			},
			validateFunc: func(context.Context, *Form) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "changer",
			repoFunc: func(m *MockRepository) {
				call := m.EXPECT().FindResetToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				call.Do(func(ctx context.Context, token string, t *Token) error {
					t.ExpiredAt = time.Now().Add(2 * time.Hour)
					return nil
				})
				call.Return(nil)
				m.EXPECT().ChangePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			validateFunc: func(context.Context, *Form) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "find",
			repoFunc: func(m *MockRepository) {
				call := m.EXPECT().FindResetToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				call.Do(func(ctx context.Context, token string, t *Token) error {
					t.ExpiredAt = time.Now().Add(2 * time.Hour)
					return nil
				})
				call.Return(nil)
				m.EXPECT().ChangePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			validateFunc: func(context.Context, *Form) error {
				return nil
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
			tt.repoFunc(repo)
			s := NewService(repo, validaterFunc(tt.validateFunc))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var f Form
			var u user.User
			err := s.Change(ctx, "", &f, &u)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

type validaterFunc func(context.Context, *Form) error

func (vf validaterFunc) Validate(ctx context.Context, f *Form) error {
	return vf(ctx, f)
}
