package authorization

import (
	"context"

	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/middleware/authentication"
)

// DefaultUserGetter the default user getter
type DefaultUserGetter struct{}

// NewDefaultUserGetter new user getter
func NewDefaultUserGetter() *DefaultUserGetter {
	return &DefaultUserGetter{}
}

// UserID implementation UserID for IUserGetter
func (u *DefaultUserGetter) UserID(ctx context.Context) string {
	authentication.FromInterceptorContext(ctx).GetUserID().String()
	return ""
}
