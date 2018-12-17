package reset

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	user "github.com/romanyx/scraper_auth/internal/user"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Reset(t *testing.T) {
	cr := func() error {
		return nil
	}

	tests := []struct {
		name       string
		repoFunc   func(m *MockRepository)
		informFunc func(context.Context, *user.User) error
		wantErr    bool
	}{
		{
			name: "ok",
			repoFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Reset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(cr, nil, nil)
			},
			informFunc: func(context.Context, *user.User) error {
				return nil
			},
		},
		{
			name: "find",
			repoFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "reset",
			repoFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Reset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "inform",
			repoFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Reset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, cr, nil)
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
			tt.repoFunc(repo)
			s := NewService(repo, informerFunc(tt.informFunc), time.Hour)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err := s.Reset(ctx, "")
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

type informerFunc func(context.Context, *user.User) error

func (f informerFunc) Inform(ctx context.Context, u *user.User) error {
	return f(ctx, u)
}
