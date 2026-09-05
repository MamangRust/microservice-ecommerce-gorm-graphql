package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func TestExtractUserIDFromClaims_String(t *testing.T) {
	claims := jwt.MapClaims{"sub": "42"}
	id, ok := extractUserIDFromClaims(claims)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if id != 42 {
		t.Fatalf("expected id=42, got %d", id)
	}
}

func TestExtractUserIDFromClaims_EmptyString(t *testing.T) {
	claims := jwt.MapClaims{"sub": ""}
	_, ok := extractUserIDFromClaims(claims)
	if ok {
		t.Fatal("expected ok=false for empty sub")
	}
}

func TestExtractUserIDFromClaims_NonNumericString(t *testing.T) {
	claims := jwt.MapClaims{"sub": "abc"}
	_, ok := extractUserIDFromClaims(claims)
	if ok {
		t.Fatal("expected ok=false for non-numeric sub")
	}
}

func TestExtractUserIDFromClaims_Float64(t *testing.T) {
	claims := jwt.MapClaims{"sub": float64(7)}
	id, ok := extractUserIDFromClaims(claims)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if id != 7 {
		t.Fatalf("expected id=7, got %d", id)
	}
}

func TestExtractUserIDFromClaims_Int(t *testing.T) {
	claims := jwt.MapClaims{"sub": 9}
	id, ok := extractUserIDFromClaims(claims)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if id != 9 {
		t.Fatalf("expected id=9, got %d", id)
	}
}

func TestExtractUserIDFromClaims_Missing(t *testing.T) {
	claims := jwt.MapClaims{"other": "value"}
	_, ok := extractUserIDFromClaims(claims)
	if ok {
		t.Fatal("expected ok=false when sub missing")
	}
}

func TestExtractUserIDFromClaims_UnsupportedType(t *testing.T) {
	claims := jwt.MapClaims{"sub": []string{"a"}}
	_, ok := extractUserIDFromClaims(claims)
	if ok {
		t.Fatal("expected ok=false for unsupported type")
	}
}

func TestSkipAuth_Whitelisted(t *testing.T) {
	paths := []string{
		"/api/auth/login",
		"/api/auth/register",
		"/api/auth/hello",
		"/api/auth/verify-code",
		"/api/auth/verify-code/anything",
		"/docs/",
		"/docs",
		"/swagger",
		"/swagger/index.html",
		"/metrics",
	}

	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, p, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if !skipAuth(c) {
				t.Fatalf("expected skipAuth=true for path %q", p)
			}
		})
	}
}

func TestSkipAuth_NonWhitelisted(t *testing.T) {
	paths := []string{
		"/api/products",
		"/api/merchants",
		"/swagge",
		"/api/auth/login-extra",
	}

	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, p, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if skipAuth(c) {
				t.Fatalf("expected skipAuth=false for path %q", p)
			}
		})
	}
}
