package verify

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	user "github.com/romanyx/scraper_auth/internal/user"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Verify(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(m *MockRepository)
		wantErr  bool
	}{
		{
			name: "ok",
			mockFunc: func(m *MockRepository) {
				m.EXPECT().FindByToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "find by token",
			mockFunc: func(m *MockRepository) {
				m.EXPECT().FindByToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "verify",
			mockFunc: func(m *MockRepository) {
				m.EXPECT().FindByToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "find",
			mockFunc: func(m *MockRepository) {
				m.EXPECT().FindByToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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
			tt.mockFunc(repo)
			s := NewService(repo)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var usr user.User
			err := s.Verify(ctx, "", &usr)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
