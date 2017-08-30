package session

import "net/http"

type SessionMiddleware struct {
	manager *Manager
	next    http.Handler
}

func (m *SessionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := ToCtx(r.Context(), m.manager.Session(w, r))
	m.next.ServeHTTP(w, r.WithContext(ctx))
}

func Middleware(manager *Manager) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &SessionMiddleware{
			manager: manager,
			next:    h,
		}
	}
}
