package apperror_test

import (
	"errors"
	"testing"

	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
	"gorm.io/gorm"
)

func TestFromGorm_NotFound(t *testing.T) {
	got := apperror.FromGorm(gorm.ErrRecordNotFound, "task not found")
	var ae *apperror.AppError
	if !errors.As(got, &ae) {
		t.Fatalf("want AppError, got %T", got)
	}
	if ae.Status != 404 {
		t.Errorf("status = %d, want 404", ae.Status)
	}
	if ae.Message != "task not found" {
		t.Errorf("message = %q", ae.Message)
	}
}

func TestFromGorm_Passthrough(t *testing.T) {
	other := errors.New("connection refused")
	got := apperror.FromGorm(other, "ignored message")
	if !errors.Is(got, other) {
		t.Errorf("non-NotFound error should pass through, got %v", got)
	}
}

func TestAppError_UnwrapsCause(t *testing.T) {
	cause := errors.New("boom")
	wrapped := apperror.Wrap(cause, 500, "internal")

	if !errors.Is(wrapped, cause) {
		t.Error("errors.Is should find the wrapped cause")
	}
}

func TestAs_Convenience(t *testing.T) {
	t.Run("returns nil ok for plain error", func(t *testing.T) {
		_, ok := apperror.As(errors.New("not an AppError"))
		if ok {
			t.Error("plain error should return ok=false")
		}
	})

	t.Run("unwraps wrapped AppError", func(t *testing.T) {
		ae := apperror.NotFound("nope")
		got, ok := apperror.As(ae)
		if !ok || got.Status != 404 {
			t.Errorf("want 404, got %v ok=%v", got, ok)
		}
	})
}
