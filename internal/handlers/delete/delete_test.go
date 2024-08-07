package delete

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"rest_api_app/internal/storage"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

type MockURLDeleter struct {
	DeleteURLFunc func(alias string) error
}

func (m *MockURLDeleter) DeleteURL(alias string) error {
	return m.DeleteURLFunc(alias)
}

func TestDeleteURL(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	urlDeleter := &MockURLDeleter{
		DeleteURLFunc: func(alias string) error {
			return nil
		},
	}

	handler := New(log, urlDeleter)

	req, err := http.NewRequest("DELETE", "/delete?alias=example", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqID := middleware.GetReqID(req.Context())
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, reqID))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"status":"OK"}`, rr.Body.String())
}

func TestDeleteURLNotFound(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	urlDeleter := &MockURLDeleter{
		DeleteURLFunc: func(alias string) error {
			return storage.ErrUrlNotFound
		},
	}

	handler := New(log, urlDeleter)

	req, err := http.NewRequest("DELETE", "/delete?alias=example", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqID := middleware.GetReqID(req.Context())
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, reqID))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{"status":"Error","error":"falied to delete url"}`, rr.Body.String())
}

func TestDeleteURLNoAlias(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	urlDeleter := &MockURLDeleter{
		DeleteURLFunc: func(alias string) error {
			return nil
		},
	}

	handler := New(log, urlDeleter)

	req, err := http.NewRequest("DELETE", "/delete", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqID := middleware.GetReqID(req.Context())
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, reqID))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"status":"Error", "error":"alias is required"}`, rr.Body.String())
}
