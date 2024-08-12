package redirect

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"rest_api_app/internal/storage"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

type MockURLGetter struct {
	GetURLFunc func(alias string) (string, error)
}

func (m *MockURLGetter) GetURL(alias string) (string, error) {
	return m.GetURLFunc(alias)
}
func TestRedirectURL(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	urlGetter := &MockURLGetter{
		GetURLFunc: func(alias string) (string, error) {
			return "http://example.com", nil
		},
	}

	handler := New(log, urlGetter)

	r := chi.NewRouter()
	r.Get("/{alias}", handler)

	req, err := http.NewRequest("GET", "/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqID := middleware.GetReqID(req.Context())
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, reqID))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, "http://example.com", rr.Header().Get("Location"))
}

func TestRedirectURLNotFound(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	urlGetter := &MockURLGetter{
		GetURLFunc: func(alias string) (string, error) {
			return "", storage.ErrUrlNotFound
		},
	}

	handler := New(log, urlGetter)

	r := chi.NewRouter()
	r.Get("/{alias}", handler)

	req, err := http.NewRequest("GET", "/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqID := middleware.GetReqID(req.Context())
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, reqID))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.JSONEq(t, `{"status":"Error","error":"not found"}`, rr.Body.String())
}
