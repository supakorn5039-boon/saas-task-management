package service

import (
	"errors"

	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{database.DB}
}

func (s *UserService) GetUserById(id uint) (*model.UserDto, error) {
	var user model.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("user not found")
		}
		return nil, apperror.Wrap(err, 500, "failed to load user")
	}
	return user.ToDto(), nil
}

// ChangePassword verifies the current password, then sets a new one.
// Both the controller (binding) and Register share the min=8 rule.
func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	var user model.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("user not found")
		}
		return apperror.Wrap(err, 500, "failed to load user")
	}

	if !security.VerifyPassword(user.Password, currentPassword) {
		return apperror.Unauthorized("current password is incorrect")
	}

	hashed, err := security.HashPassword(newPassword)
	if err != nil {
		return apperror.Wrap(err, 500, "failed to hash password")
	}
	if err := s.db.Model(&user).Update("password", hashed).Error; err != nil {
		return apperror.Wrap(err, 500, "failed to update password")
	}
	return nil
}

// ----- Admin-only operations -----

var allowedUserSortColumns = map[string]string{
	"created_at": "created_at",
	"email":      "email",
	"role":       "role",
	"status":     "status",
}

type ListUsersOptions struct {
	Page    int
	PerPage int
	Search  string
	Sort    string
	Order   string
}

func (s *UserService) ListUsers(opts ListUsersOptions) (*model.UserListResponse, error) {
	q := s.db.Model(&model.User{})
	if opts.Search != "" {
		like := "%" + opts.Search + "%"
		q = q.Where("email ILIKE ?", like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "failed to count users")
	}

	sortCol, ok := allowedUserSortColumns[opts.Sort]
	if !ok {
		sortCol = "created_at"
	}
	order := "desc"
	if opts.Order == "asc" {
		order = "asc"
	}

	var users []model.User
	err := q.Order(sortCol + " " + order).
		Limit(opts.PerPage).
		Offset((opts.Page - 1) * opts.PerPage).
		Find(&users).Error
	if err != nil {
		return nil, apperror.Wrap(err, 500, "failed to list users")
	}

	dtos := make([]*model.UserDto, len(users))
	for i, u := range users {
		dtos[i] = u.ToDto()
	}

	return &model.UserListResponse{
		Data: dtos,
		Meta: model.UserListMeta{
			Page:    opts.Page,
			PerPage: opts.PerPage,
			Total:   total,
		},
	}, nil
}

type AdminUpdateUserInput struct {
	Role   *string
	Status *int
}

// AdminUpdateUser changes role and/or status for any user. The actorID is the
// admin performing the action — they cannot demote themselves out of admin
// or deactivate their own account (would lock themselves out).
func (s *UserService) AdminUpdateUser(actorID, targetID uint, in AdminUpdateUserInput) (*model.UserDto, error) {
	if in.Role == nil && in.Status == nil {
		return nil, apperror.BadRequest("no fields to update")
	}
	if in.Role != nil && !model.IsValidRole(*in.Role) {
		return nil, apperror.BadRequest("invalid role")
	}

	var user model.User
	if err := s.db.First(&user, "id = ?", targetID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("user not found")
		}
		return nil, apperror.Wrap(err, 500, "failed to load user")
	}

	if actorID == targetID {
		if in.Role != nil && *in.Role != "admin" {
			return nil, apperror.BadRequest("you cannot demote yourself")
		}
		if in.Status != nil && *in.Status != 1 {
			return nil, apperror.BadRequest("you cannot deactivate yourself")
		}
	}

	if in.Role != nil {
		user.Role = *in.Role
	}
	if in.Status != nil {
		user.Status = *in.Status
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "failed to update user")
	}
	return user.ToDto(), nil
}

// ListAssignable returns the lightweight user list used to populate the
// task-assignee dropdown. Active users only — deactivated accounts can't be
// assigned new work. No auth check here; the controller mounts this behind
// Protected so any authenticated user can call it.
func (s *UserService) ListAssignable() ([]*model.UserDto, error) {
	var users []model.User
	err := s.db.Where("status = ?", 1).
		Order("email asc").
		Find(&users).Error
	if err != nil {
		return nil, apperror.Wrap(err, 500, "failed to list users")
	}
	dtos := make([]*model.UserDto, len(users))
	for i, u := range users {
		dtos[i] = u.ToDto()
	}
	return dtos, nil
}

// AdminDeleteUser soft-deletes a user (gorm.Model has DeletedAt). Self-delete
// is rejected to avoid locking the actor out.
func (s *UserService) AdminDeleteUser(actorID, targetID uint) error {
	if actorID == targetID {
		return apperror.BadRequest("you cannot delete your own account")
	}
	result := s.db.Delete(&model.User{}, "id = ?", targetID)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, "failed to delete user")
	}
	if result.RowsAffected == 0 {
		return apperror.NotFound("user not found")
	}
	return nil
}
