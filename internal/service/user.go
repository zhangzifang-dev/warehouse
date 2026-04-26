package service

import (
	"context"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/pkg/password"
)

type UserFullRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context, page, pageSize int) ([]model.User, int, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]model.Role, error)
	AssignRoles(ctx context.Context, userID int64, roleIDs []int64) error
}

type UserService struct {
	userRepo UserFullRepository
}

func NewUserService(userRepo UserFullRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type CreateUserInput struct {
	Username string
	Password string
	Status   int
}

func (s *UserService) Create(ctx context.Context, input *CreateUserInput) (*model.User, error) {
	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to hash password")
	}

	user := &model.User{
		Username: input.Username,
		Password: hashedPassword,
		Status:   input.Status,
	}

	if user.Status == 0 {
		user.Status = model.UserStatusActive
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		if isDuplicateEntry(err) {
			return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "username already exists")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create user")
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
	}
	return user, nil
}

type ListUsersResult struct {
	Users []model.User
	Total int
}

func (s *UserService) List(ctx context.Context, page, pageSize int) (*ListUsersResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	users, total, err := s.userRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list users")
	}

	return &ListUsersResult{
		Users: users,
		Total: total,
	}, nil
}

type UpdateUserInput struct {
	Status *int
}

func (s *UserService) Update(ctx context.Context, id int64, input *UpdateUserInput) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
	}

	if input.Status != nil {
		user.Status = *input.Status
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update user")
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
	}

	err = s.userRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete user")
	}

	return nil
}

func (s *UserService) GetUserRoles(ctx context.Context, userID int64) ([]model.Role, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
	}

	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get user roles")
	}

	return roles, nil
}

func (s *UserService) AssignRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
	}

	err = s.userRepo.AssignRoles(ctx, userID, roleIDs)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to assign roles")
	}

	return nil
}

func isDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "duplicate key" || err.Error() == "ERROR #1062"
}
