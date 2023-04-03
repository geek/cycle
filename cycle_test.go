package cycle_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geek/cycle"
	"github.com/geek/herrors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func handler(w http.ResponseWriter, r *http.Request) {}

func TestUseWithMux(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")

	c := cycle.New(r)

	called := 0
	c.OnRequest(func(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
		called++
		w.WriteHeader(201)
		return r, nil
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

func TestNotFoundUrl(t *testing.T) {
	r := mux.NewRouter()

	c := cycle.New(r)

	called := 0
	c.OnRequest(func(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
		called++
		return r, nil
	})

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	defer req.Body.Close()
	r.ServeHTTP(rw, req)

	assert.Equal(t, 1, called)
	assert.Equal(t, 404, rw.Result().StatusCode)
}

func TestUpdateContext(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")

	c := cycle.New(r)

	ctxKey := "rcalled"
	called := 0
	c.OnRequest(func(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
		called++
		rcalled := 1
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKey, rcalled)

		return r.WithContext(ctx), nil
	})

	c.OnRequest(func(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
		called++
		ctx := r.Context()
		rcalled := ctx.Value(ctxKey).(int)
		rcalled++
		ctx = context.WithValue(ctx, ctxKey, rcalled)

		w.WriteHeader(200 + rcalled)

		return r.WithContext(ctx), nil
	})

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	defer req.Body.Close()
	r.ServeHTTP(rw, req)

	assert.Equal(t, 2, called)
	assert.Equal(t, 200+called, rw.Result().StatusCode)
}

func TestAuthIsPopulatedForHandler(t *testing.T) {
	r := mux.NewRouter()
	c := cycle.New(r)

	ctxKey := "auth"
	user := "user1"
	c.OnAuth(func(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKey, user)

		return r.WithContext(ctx), nil
	})

	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ruser := ctx.Value(ctxKey).(string)
		assert.Equal(t, user, ruser)
	}

	r.HandleFunc("/", handler).Methods("GET")

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	defer req.Body.Close()
	r.ServeHTTP(rw, req)

	assert.Equal(t, 200, rw.Result().StatusCode)
}

func TestHTTPError(t *testing.T) {
	r := mux.NewRouter()

	called := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		called++
	}

	r.HandleFunc("/", handler).Methods("GET")

	c := cycle.New(r)

	c.OnRequest(func(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
		return nil, herrors.ErrBadRequest
	})

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	defer req.Body.Close()
	r.ServeHTTP(rw, req)

	assert.Equal(t, 0, called)
	assert.Equal(t, 400, rw.Result().StatusCode)
}

func TestUnkownError(t *testing.T) {
	r := mux.NewRouter()

	called := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		called++
	}

	r.HandleFunc("/", handler).Methods("GET")

	c := cycle.New(r)

	c.OnRequest(func(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
		return nil, errors.New("my error")
	})

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	defer req.Body.Close()
	r.ServeHTTP(rw, req)

	assert.Equal(t, 0, called)
	assert.Equal(t, 500, rw.Result().StatusCode)
}
