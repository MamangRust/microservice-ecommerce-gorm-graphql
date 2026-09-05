package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type stubRoleCache struct {
	roles map[string][]string
}

func (s *stubRoleCache) GetRoleCache(_ context.Context, userID string) ([]string, bool) {
	roles, ok := s.roles[userID]
	return roles, ok
}

func (s *stubRoleCache) SetRoleCache(_ context.Context, userID string, roles []string) {
	s.roles[userID] = roles
}

func nopLogger() logger.LoggerInterface {
	return &logger.Logger{Log: zap.NewNop()}
}

func TestExtractUserID_Float64(t *testing.T) {
	v := &RoleValidator{}
	id, err := v.extractUserID(float64(123))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 123 {
		t.Fatalf("expected 123, got %d", id)
	}
}

func TestExtractUserID_Int(t *testing.T) {
	v := &RoleValidator{}
	id, err := v.extractUserID(55)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 55 {
		t.Fatalf("expected 55, got %d", id)
	}
}

func TestExtractUserID_String(t *testing.T) {
	v := &RoleValidator{}
	id, err := v.extractUserID("77")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 77 {
		t.Fatalf("expected 77, got %d", id)
	}
}

func TestExtractUserID_InvalidString(t *testing.T) {
	v := &RoleValidator{}
	_, err := v.extractUserID("not-a-number")
	if err == nil {
		t.Fatal("expected error for non-numeric string")
	}
}

func TestExtractUserID_UnknownType(t *testing.T) {
	v := &RoleValidator{}
	_, err := v.extractUserID([]string{"a"})
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestRoleValidatorMiddleware_CacheHit(t *testing.T) {
	v := &RoleValidator{
		logger: nopLogger(),
		cache: &stubRoleCache{
			roles: map[string][]string{"123": {"admin", "merchant"}},
		},
		timeout: 2 * time.Second,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", 123)

	nextCalled := false
	handler := v.Middleware()(func(c echo.Context) error {
		nextCalled = true
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}

	roleNames, ok := c.Get("role_names").([]string)
	if !ok {
		t.Fatal("expected role_names set in context")
	}
	if len(roleNames) != 2 || roleNames[0] != "admin" {
		t.Fatalf("unexpected role_names: %v", roleNames)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestRoleValidatorMiddleware_MissingUserID(t *testing.T) {
	v := &RoleValidator{
		logger: nopLogger(),
		cache:  &stubRoleCache{roles: map[string][]string{}},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	nextCalled := false
	handler := v.Middleware()(func(c echo.Context) error {
		nextCalled = true
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	he, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}
	if he.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", he.Code)
	}
	if nextCalled {
		t.Fatal("expected next handler NOT to be called")
	}
}

func TestRoleValidatorMiddleware_InvalidUserID(t *testing.T) {
	v := &RoleValidator{
		logger: nopLogger(),
		cache:  &stubRoleCache{roles: map[string][]string{}},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "not-a-number")

	nextCalled := false
	handler := v.Middleware()(func(c echo.Context) error {
		nextCalled = true
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	he, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}
	if he.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", he.Code)
	}
	if nextCalled {
		t.Fatal("expected next handler NOT to be called")
	}
}
