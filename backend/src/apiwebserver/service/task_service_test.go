package service_test

import (
	"errors"
	"testing"

	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"github.com/supakorn5039-boon/saas-task-backend/src/testhelpers"
)

func registerForTest(t *testing.T, email string) uint {
	t.Helper()
	auth := service.NewAuthService()
	user, err := auth.Register(email, "password123")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	return user.Id
}

func TestTaskService_ListTasks(t *testing.T) {
	testhelpers.SetupTestDB(t)
	userID := registerForTest(t, "list@example.com")
	otherID := registerForTest(t, "other@example.com")

	tasks := service.NewTaskService()
	if _, err := tasks.CreateTask(userID, "Buy milk", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := tasks.CreateTask(userID, "Write tests", "for the service layer"); err != nil {
		t.Fatal(err)
	}
	// Different owner — must not appear in our list.
	if _, err := tasks.CreateTask(otherID, "Not yours", ""); err != nil {
		t.Fatal(err)
	}

	t.Run("returns only own tasks with counts", func(t *testing.T) {
		res, err := tasks.ListTasks(service.ListTasksOptions{
			UserID:  userID,
			Page:    1,
			PerPage: 10,
		})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if res.Meta.Total != 2 {
			t.Errorf("total = %d, want 2", res.Meta.Total)
		}
		if res.Counts.Todo != 2 || res.Counts.All != 2 {
			t.Errorf("counts = %+v, want all=2 todo=2", res.Counts)
		}
	})

	t.Run("search is case-insensitive across title and description", func(t *testing.T) {
		res, err := tasks.ListTasks(service.ListTasksOptions{
			UserID:  userID,
			Page:    1,
			PerPage: 10,
			Search:  "TESTS",
		})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(res.Data) != 1 || res.Data[0].Title != "Write tests" {
			t.Errorf("got %+v, want [Write tests]", res.Data)
		}
	})

	t.Run("invalid sort column falls back to created_at", func(t *testing.T) {
		// Should not error; sortable column whitelist protects against SQL injection
		// and silently corrects bad input.
		_, err := tasks.ListTasks(service.ListTasksOptions{
			UserID:  userID,
			Page:    1,
			PerPage: 10,
			Sort:    "name; DROP TABLE tasks; --",
		})
		if err != nil {
			t.Fatalf("expected silent fallback, got %v", err)
		}
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	testhelpers.SetupTestDB(t)
	userID := registerForTest(t, "owner@example.com")
	otherID := registerForTest(t, "intruder@example.com")

	tasks := service.NewTaskService()
	created, err := tasks.CreateTask(userID, "Old", "old desc")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("partial patch only changes provided fields", func(t *testing.T) {
		title := "New title"
		updated, err := tasks.UpdateTask(userID, created.ID, service.UpdateTaskInput{
			Title: &title,
		})
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		if updated.Title != "New title" {
			t.Errorf("title = %q", updated.Title)
		}
		if updated.Description != "old desc" {
			t.Errorf("description should be untouched, got %q", updated.Description)
		}
	})

	t.Run("status flip", func(t *testing.T) {
		s := model.TaskStatusInProgress
		updated, err := tasks.UpdateTask(userID, created.ID, service.UpdateTaskInput{Status: &s})
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		if updated.Status != model.TaskStatusInProgress {
			t.Errorf("status = %q", updated.Status)
		}
	})

	t.Run("foreign user gets 404 not 403", func(t *testing.T) {
		// Privacy: don't reveal that the task exists. Stick to NotFound.
		title := "hijack"
		_, err := tasks.UpdateTask(otherID, created.ID, service.UpdateTaskInput{Title: &title})
		var ae *apperror.AppError
		if !errors.As(err, &ae) {
			t.Fatalf("want AppError, got %T", err)
		}
		if ae.Status != 404 {
			t.Errorf("status = %d, want 404", ae.Status)
		}
	})
}

func TestTaskService_DeleteTask(t *testing.T) {
	testhelpers.SetupTestDB(t)
	userID := registerForTest(t, "del@example.com")
	tasks := service.NewTaskService()
	created, _ := tasks.CreateTask(userID, "to delete", "")

	if err := tasks.DeleteTask(userID, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	// Second delete: row already soft-deleted → not found.
	err := tasks.DeleteTask(userID, created.ID)
	var ae *apperror.AppError
	if !errors.As(err, &ae) || ae.Status != 404 {
		t.Errorf("second delete: want 404, got %v", err)
	}
}
