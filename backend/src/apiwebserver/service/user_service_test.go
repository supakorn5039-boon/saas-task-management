package service_test

import (
	"errors"
	"testing"

	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
	"github.com/supakorn5039-boon/saas-task-backend/src/testhelpers"
)

func TestUserService_ChangePassword(t *testing.T) {
	testhelpers.SetupTestDB(t)
	userID := registerForTest(t, "pw@example.com")

	users := service.NewUserService()
	auth := service.NewAuthService()

	t.Run("wrong current password rejected", func(t *testing.T) {
		err := users.ChangePassword(userID, "wrong", "another-strong-pass")
		var ae *apperror.AppError
		if !errors.As(err, &ae) || ae.Status != 401 {
			t.Fatalf("want 401, got %v", err)
		}
	})

	t.Run("happy path lets new password log in", func(t *testing.T) {
		if err := users.ChangePassword(userID, "password123", "another-strong-pass"); err != nil {
			t.Fatalf("change: %v", err)
		}
		if _, err := auth.Login("pw@example.com", "another-strong-pass"); err != nil {
			t.Errorf("expected login to succeed with new password, got %v", err)
		}
		if _, err := auth.Login("pw@example.com", "password123"); err == nil {
			t.Error("expected old password to fail after change")
		}
	})
}

func TestUserService_AdminOperations(t *testing.T) {
	testhelpers.SetupTestDB(t)

	adminID := registerForTest(t, "admin@example.com")
	memberID := registerForTest(t, "member@example.com")

	users := service.NewUserService()

	t.Run("admin can promote a member", func(t *testing.T) {
		role := "manager"
		updated, err := users.AdminUpdateUser(adminID, memberID, service.AdminUpdateUserInput{Role: &role})
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		if updated.Role != "manager" {
			t.Errorf("role = %q", updated.Role)
		}
	})

	t.Run("invalid role rejected with 400", func(t *testing.T) {
		bad := "superuser"
		_, err := users.AdminUpdateUser(adminID, memberID, service.AdminUpdateUserInput{Role: &bad})
		var ae *apperror.AppError
		if !errors.As(err, &ae) || ae.Status != 400 {
			t.Errorf("want 400, got %v", err)
		}
	})

	t.Run("admin cannot demote themselves", func(t *testing.T) {
		// The actor here is acting on their own row — must be blocked.
		role := "user"
		_, err := users.AdminUpdateUser(adminID, adminID, service.AdminUpdateUserInput{Role: &role})
		var ae *apperror.AppError
		if !errors.As(err, &ae) || ae.Status != 400 {
			t.Errorf("want 400 self-demote, got %v", err)
		}
	})

	t.Run("admin cannot deactivate themselves", func(t *testing.T) {
		inactive := 0
		_, err := users.AdminUpdateUser(adminID, adminID, service.AdminUpdateUserInput{Status: &inactive})
		var ae *apperror.AppError
		if !errors.As(err, &ae) || ae.Status != 400 {
			t.Errorf("want 400 self-deactivate, got %v", err)
		}
	})

	t.Run("admin cannot delete themselves", func(t *testing.T) {
		err := users.AdminDeleteUser(adminID, adminID)
		var ae *apperror.AppError
		if !errors.As(err, &ae) || ae.Status != 400 {
			t.Errorf("want 400 self-delete, got %v", err)
		}
	})

	t.Run("admin can delete a member", func(t *testing.T) {
		if err := users.AdminDeleteUser(adminID, memberID); err != nil {
			t.Fatalf("delete: %v", err)
		}
		// Already deleted → not found.
		err := users.AdminDeleteUser(adminID, memberID)
		var ae *apperror.AppError
		if !errors.As(err, &ae) || ae.Status != 404 {
			t.Errorf("want 404 on second delete, got %v", err)
		}
	})
}

func TestUserService_ListUsers(t *testing.T) {
	testhelpers.SetupTestDB(t)
	registerForTest(t, "alpha@example.com")
	registerForTest(t, "beta@example.com")
	registerForTest(t, "gamma@example.com")

	users := service.NewUserService()

	t.Run("paginates and counts", func(t *testing.T) {
		res, err := users.ListUsers(service.ListUsersOptions{Page: 1, PerPage: 2})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if res.Meta.Total != 3 {
			t.Errorf("total = %d, want 3", res.Meta.Total)
		}
		if len(res.Data) != 2 {
			t.Errorf("len(data) = %d, want 2", len(res.Data))
		}
	})

	t.Run("search filters by email", func(t *testing.T) {
		res, err := users.ListUsers(service.ListUsersOptions{Page: 1, PerPage: 10, Search: "BETA"})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if res.Meta.Total != 1 {
			t.Errorf("total = %d, want 1", res.Meta.Total)
		}
	})
}
