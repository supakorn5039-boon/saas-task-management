package service_test

import (
	"errors"
	"testing"

	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
	"github.com/supakorn5039-boon/saas-task-backend/src/testhelpers"
)

func TestAuthService_Register(t *testing.T) {
	testhelpers.SetupTestDB(t)
	svc := service.NewAuthService()

	t.Run("happy path", func(t *testing.T) {
		user, err := svc.Register("alice@example.com", "password123")
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if user.Email != "alice@example.com" {
			t.Errorf("email = %q, want alice@example.com", user.Email)
		}
		if user.Role != "user" {
			t.Errorf("role = %q, want user", user.Role)
		}
	})

	t.Run("duplicate email returns 409 AppError", func(t *testing.T) {
		_, err := svc.Register("alice@example.com", "password123")
		var ae *apperror.AppError
		if !errors.As(err, &ae) {
			t.Fatalf("want AppError, got %T", err)
		}
		if ae.Status != 409 {
			t.Errorf("status = %d, want 409", ae.Status)
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	testhelpers.SetupTestDB(t)
	svc := service.NewAuthService()

	if _, err := svc.Register("bob@example.com", "password123"); err != nil {
		t.Fatalf("setup: %v", err)
	}

	t.Run("happy path", func(t *testing.T) {
		user, err := svc.Login("bob@example.com", "password123")
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if user.Id == 0 {
			t.Error("expected non-zero user id")
		}
	})

	t.Run("wrong password returns 401 with generic message", func(t *testing.T) {
		_, err := svc.Login("bob@example.com", "wrong")
		var ae *apperror.AppError
		if !errors.As(err, &ae) {
			t.Fatalf("want AppError, got %T", err)
		}
		if ae.Status != 401 {
			t.Errorf("status = %d, want 401", ae.Status)
		}
		// Don't leak whether the email exists — message is the same as for unknown email.
		if ae.Message != "invalid email or password" {
			t.Errorf("message = %q, want generic", ae.Message)
		}
	})

	t.Run("unknown email returns the same 401", func(t *testing.T) {
		_, err := svc.Login("nobody@example.com", "whatever")
		var ae *apperror.AppError
		if !errors.As(err, &ae) || ae.Status != 401 {
			t.Fatalf("want 401 AppError, got %v", err)
		}
		if ae.Message != "invalid email or password" {
			t.Errorf("message = %q, want generic", ae.Message)
		}
	})
}
