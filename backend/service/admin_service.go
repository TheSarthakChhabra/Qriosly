package service

import (
	"context"
	"quiz-backend/apperror"
	"quiz-backend/model"
)

type AdminService struct {
	userRepo UserRepo
}

func NewAdminService(userRepo UserRepo) *AdminService {
	return &AdminService{userRepo: userRepo}
}

func (s *AdminService) ListUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, apperror.Internal("USERS_FETCH_FAILED", "failed to fetch users")
	}
	return users, nil
}

// UpdateUserRole is the ONLY place role changes happen post-registration.
func (s *AdminService) UpdateUserRole(ctx context.Context, requesterID, targetUserID, newRole string) (model.User, error) {
	if requesterID == targetUserID {
		return model.User{}, apperror.BadRequest("CANNOT_MODIFY_OWN_ROLE", "you cannot change your own role")
	}

	if newRole != model.RoleStudent && newRole != model.RoleTeacher && newRole != model.RoleAdmin {
		return model.User{}, apperror.BadRequest("INVALID_ROLE", "invalid role")
	}

	user, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return model.User{}, apperror.NotFound("USER_NOT_FOUND", "user not found")
	}

	if err := s.userRepo.UpdateRole(ctx, targetUserID, newRole); err != nil {
		return model.User{}, apperror.Internal("ROLE_UPDATE_FAILED", "failed to update role")
	}

	user.Role = newRole
	return user, nil
}
