package authz

import "context"

type DefaultUserGetter struct {
}

func NewDefaultUserGetter() *DefaultUserGetter {
	return &DefaultUserGetter{}
}

func (u *DefaultUserGetter) UserID(ctx context.Context) string {
	// authorization.FromInterceptorContext(ctx).GetUserID().String()
	return ""
}
