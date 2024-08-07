package save

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

type MockURLSaver struct {
	SaveURLFunc func(urlToSave string, alias string) (int64, error)
}

func (m *MockURLSaver) SaveURL(urlToSave string, alias string) (int64, error) {
	return m.SaveURLFunc(urlToSave, alias)
}

// тест на то что создается все как надо
func TestSaveURLWithAlias(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	urlSaver := &MockURLSaver{
		SaveURLFunc: func(urlToSave string, alias string) (int64, error) {
			return 1, nil
		},
	}

	handler := New(log, urlSaver)

	reqBody := `{"url": "http://example.com", "alias": "example"}`
	req, err := http.NewRequest("POST", "/save", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	reqID := middleware.GetReqID(req.Context())
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, reqID))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"status":"OK","alias":"example"}`, rr.Body.String())
}

// Тест на то что создается случайное название для alias если оно не указано
func TestSaveUrlWithoutAlias(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	urlSaver := &MockURLSaver{
		SaveURLFunc: func(urlToSave string, alias string) (int64, error) {
			return 1, nil
		},
	}

	handler := New(log, urlSaver)

	reqBody := `{"url": "http://example.com"}`
	req, err := http.NewRequest("POST", "/save", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	reqID := middleware.GetReqID(req.Context())
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, reqID))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responce map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responce)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "OK", responce["status"])
	assert.NotEmpty(t, responce["alias"])
}

func TestSaveUrlWithInvalidJson(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	urlSaver := &MockURLSaver{
		SaveURLFunc: func(urlToSave string, alias string) (int64, error) {
			return 1, nil
		},
	}

	handler := New(log, urlSaver)

	reqBody := `{"url": "http://example.com", "alias": "example"`
	req, err := http.NewRequest("POST", "/save", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	reqID := middleware.GetReqID(req.Context())
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, reqID))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"status":"Error","error":"falied to decode request"}`, rr.Body.String())
}
