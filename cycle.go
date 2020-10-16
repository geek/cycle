package cycle

import (
	"net/http"
)

func New() *Cycle {
	return &Cycle{}
}

type Cycle struct {
	onRequest []http.HandlerFunc
}

func (c *Cycle) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, onRequest := range c.onRequest {
			onRequest(w, r)
		}

		h.ServeHTTP(w, r)
	})
}

func (c *Cycle) OnRequest(h http.HandlerFunc) {
	c.onRequest = append(c.onRequest, h)
}
