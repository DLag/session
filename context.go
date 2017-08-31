package session

import "context"

type sessionContextKey struct {}

func FromCtx(ctx context.Context) (s *Session) {
	s, _ = ctx.Value(sessionContextKey{}).(*Session)
	return
}

func ToCtx(ctx context.Context, s *Session) (context.Context){
	return context.WithValue(ctx, sessionContextKey{},s)
}