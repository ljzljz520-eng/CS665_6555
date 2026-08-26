package httpapi

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"scriptstudio/script-backend/internal/query"
	"scriptstudio/script-backend/internal/service"
	"scriptstudio/script-backend/internal/store"
	"strings"
	"testing"
)

func TestHTTPCreateAndList(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "http.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	api := New(service.New(repo), query.New(repo))
	req := httptest.NewRequest(http.MethodPost, "/api/scripts", strings.NewReader(`{"requestKey":"http-1","title":"Red Door","genre":"noir"}`))
	rec := httptest.NewRecorder()
	api.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", rec.Code, rec.Body.String())
	}
	req = httptest.NewRequest(http.MethodGet, "/api/scripts?q=door", nil)
	rec = httptest.NewRecorder()
	api.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Red Door") {
		t.Fatalf("list response: %d %s", rec.Code, rec.Body.String())
	}
}
