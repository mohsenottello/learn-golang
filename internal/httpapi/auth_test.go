package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestProtectedRoutesRequireLogin(t *testing.T) {
	router := NewRouter(Deps{Static: fstest.MapFS{}})

	for _, path := range []string{"/api/auth/me", "/api/me/resume", "/api/admin/users", "/api/admin/site"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("%s: expected 401, got %d", path, rec.Code)
		}
	}
}

func TestTokenHashIsStable(t *testing.T) {
	raw, hash := newToken()
	if hashToken(raw) != hash {
		t.Fatal("hashToken(raw) must equal the stored hash")
	}
}
