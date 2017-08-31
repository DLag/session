package session

import "net/http"

type Middleware struct {
	manager *Manager
	next    http.Handler
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := ToCtx(r.Context(), m.manager.Session(w, r))
	m.next.ServeHTTP(w, r.WithContext(ctx))
}

func NewMiddleware(manager *Manager) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Middleware{
			manager: manager,
			next:    h,
		}
	}
}
