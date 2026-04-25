package service

import (
	"context"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
	appjwt "warehouse/internal/pkg/jwt"
	"warehouse/internal/pkg/password"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type AuthService struct {
	userRepo UserRepository
	jwt      *appjwt.JWT
}

func NewAuthService(userRepo UserRepository, jwt *appjwt.JWT) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (s *AuthService) Login(ctx context.Context, username, pwd string) (string, *model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
	}

	if !password.Verify(pwd, user.Password) {
		return "", nil, apperrors.NewAppError(apperrors.CodeInvalidPassword, "invalid password")
	}

	if !user.IsActive() {
		return "", nil, apperrors.NewAppError(apperrors.CodeForbidden, "user is disabled")
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to generate token")
	}

	return token, user, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
	}

	return user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
	}

	if !user.IsActive() {
		return apperrors.NewAppError(apperrors.CodeForbidden, "user is disabled")
	}

	if !password.Verify(oldPassword, user.Password) {
		return apperrors.NewAppError(apperrors.CodeInvalidPassword, "invalid old password")
	}

	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to hash password")
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to update password")
	}

	return nil
}
