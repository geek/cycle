package cycle

import (
	"net/http"

	"github.com/gorilla/mux"
)

func New(r *mux.Router) *Cycle {
	c := &Cycle{r: r}
	r.Use(c.middleware)
	r.NotFoundHandler = c.notFoundHandler()

	return c
}

type Cycle struct {
	r          *mux.Router
	onRequest  []http.HandlerFunc
	onPreAuth  []http.HandlerFunc
	onAuth     []http.HandlerFunc
	onPostAuth []http.HandlerFunc
	onValidate []http.HandlerFunc
}

func (c *Cycle) middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, onRequest := range c.onRequest {
			onRequest(w, r)
		}

		for _, onPreAuth := range c.onPreAuth {
			onPreAuth(w, r)
		}

		for _, onAuth := range c.onAuth {
			onAuth(w, r)
		}

		for _, onPostAuth := range c.onPostAuth {
			onPostAuth(w, r)
		}

		for _, onValidate := range c.onValidate {
			onValidate(w, r)
		}

		h.ServeHTTP(w, r)
	})
}

func (c *Cycle) notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, onRequest := range c.onRequest {
			onRequest(w, r)
		}

		http.NotFoundHandler().ServeHTTP(w, r)
	})
}

func (c *Cycle) OnRequest(h http.HandlerFunc) {
	c.onRequest = append(c.onRequest, h)
}

func (c *Cycle) OnPreAuth(h http.HandlerFunc) {
	c.onPreAuth = append(c.onPreAuth, h)
}

func (c *Cycle) OnAuth(h http.HandlerFunc) {
	c.onAuth = append(c.onAuth, h)
}

func (c *Cycle) OnPostAuth(h http.HandlerFunc) {
	c.onPostAuth = append(c.onPostAuth, h)
}

func (c *Cycle) OnValidate(h http.HandlerFunc) {
	c.onValidate = append(c.onValidate, h)
}
