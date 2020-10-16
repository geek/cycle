package cycle

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func handler(w http.ResponseWriter, r *http.Request) {}

func TestUseWithMux(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")

	c := New()
	r.Use(c.Middleware)

	called := 0
	c.OnRequest(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(201)
	})

	t.Run("onRequest is handled", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		defer req.Body.Close()
		r.ServeHTTP(rw, req)

		assert.Equal(t, 1, called)
		assert.Equal(t, 201, rw.Result().StatusCode)
	})
}
