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

type HandlerFunc func(w http.ResponseWriter, r *http.Request) (*http.Request, error)

type Cycle struct {
	r          *mux.Router
	onRequest  []HandlerFunc
	onPreAuth  []HandlerFunc
	onAuth     []HandlerFunc
	onPostAuth []HandlerFunc
	onValidate []HandlerFunc
}

func (c *Cycle) handleError(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: use boom to write error
}

func (c *Cycle) middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		for _, onRequest := range c.onRequest {
			r, err = onRequest(w, r)
			if err != nil {
				c.handleError(w, r, err)
				return
			}

		}

		for _, onPreAuth := range c.onPreAuth {
			r, err = onPreAuth(w, r)
			if err != nil {
				c.handleError(w, r, err)
				return
			}
		}

		for _, onAuth := range c.onAuth {
			r, err = onAuth(w, r)
			if err != nil {
				c.handleError(w, r, err)
				return
			}
		}

		for _, onPostAuth := range c.onPostAuth {
			r, err = onPostAuth(w, r)
			if err != nil {
				c.handleError(w, r, err)
				return
			}
		}

		for _, onValidate := range c.onValidate {
			r, err = onValidate(w, r)
			if err != nil {
				c.handleError(w, r, err)
				return
			}
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

func (c *Cycle) OnRequest(h HandlerFunc) {
	c.onRequest = append(c.onRequest, h)
}

func (c *Cycle) OnPreAuth(h HandlerFunc) {
	c.onPreAuth = append(c.onPreAuth, h)
}

func (c *Cycle) OnAuth(h HandlerFunc) {
	c.onAuth = append(c.onAuth, h)
}

func (c *Cycle) OnPostAuth(h HandlerFunc) {
	c.onPostAuth = append(c.onPostAuth, h)
}

func (c *Cycle) OnValidate(h HandlerFunc) {
	c.onValidate = append(c.onValidate, h)
}
