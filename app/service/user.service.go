package service

import (
	"context"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/package/errors"
	"sondth-test_soa/utils"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type userService struct {
	postgresRepo repository.RepositoryCollections
	helper       helper.HelperCollections
}

func NewUserService(
	postgresRepo repository.RepositoryCollections,
	helper helper.HelperCollections,
) IUserService {
	return &userService{
		helper:       helper,
		postgresRepo: postgresRepo,
	}
}

func (s *userService) Register(
	ctx context.Context,
	req *model.UserRegisterRequest,
) (*model.UserRegisterResponse, error) {
	// Check username
	existedUser, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
		Username: &req.Username,
		Filter: repository.Filter{
			Fields: []string{"id"},
		},
	})
	if existedUser != nil || err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New(errors.ErrCodeUserExisted)
	}

	user := entity.NewUser()
	user.Username = req.Username
	user.Password = req.Password
	user.Fullname = req.Fullname
	user.Role = entity.ROLE_USER
	if err := s.postgresRepo.UserRepo.Create(ctx, nil, user); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.UserRegisterResponse{}, nil
}

func (s *userService) Login(
	ctx context.Context,
	req *model.UserLoginRequest,
) (*model.UserLoginResponse, error) {
	// Find user by username
	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
		Username: &req.Username,
		Filter: repository.Filter{
			Fields: []string{"id", "password"},
		},
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New(errors.ErrCodeIncorrectPassword)
	}

	// Generate tokens
	accessToken, err := s.helper.OAuthHelper.GenerateAccessToken(*user)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	refreshToken, err := s.helper.OAuthHelper.GenerateRefreshToken(*user)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) ForgetPassword(
	ctx context.Context,
	req *model.ForgetPasswordRequest,
) (*model.ForgetPasswordResponse, error) {
	// Find user by username
	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
		Username: &req.Username,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}

	// Update password
	user.Password = req.NewPassword

	// Update user
	if err := s.postgresRepo.UserRepo.Update(ctx, nil, user); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.ForgetPasswordResponse{}, nil
}

func (s *userService) ChangePassword(
	ctx context.Context,
	req *model.ChangeUserPasswordRequest,
) (*model.ChangeUserPasswordResponse, error) {
	if !s.helper.UserHelper.IsValidPermission(ctx, req.ID) {
		return nil, errors.New(errors.ErrCodeUnauthorized)
	}

	// Find user by ID
	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
		ID: &req.ID,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}

	// Check old password
	if err := user.CheckPassword(req.OldPassword); err != nil {
		return nil, errors.New(errors.ErrCodeIncorrectPassword)
	}

	// Update password
	user.Password = req.NewPassword
	if err := s.postgresRepo.UserRepo.Update(ctx, nil, user); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.ChangeUserPasswordResponse{}, nil
}

func (s *userService) UpdateUser(
	ctx context.Context,
	req *model.UpdateUserRequest,
) (*model.UpdateUserResponse, error) {
	if !s.helper.UserHelper.IsValidPermission(ctx, req.ID) {
		return nil, errors.New(errors.ErrCodeUnauthorized)
	}

	requestUser, ok := ctx.Value(string(utils.USER_CONTEXT_KEY)).(*entity.User)
	if !ok {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Find user by ID
	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
		ID: &req.ID,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}

	// Update fields if provided
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Fullname != nil {
		user.Fullname = *req.Fullname
	}
	if req.Role != nil && requestUser.IsAdmin() {
		if !s.helper.UserHelper.IsValidRole(*req.Role) {
			return nil, errors.NewCustomError(
				errors.ErrCodeValidatorFormat,
				errors.GetCustomMessage(errors.ErrCodeValidatorFormat, "Role"),
			)
		}
		user.Role = *req.Role
	}

	// Update user
	if err := s.postgresRepo.UserRepo.Update(ctx, nil, user); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.UpdateUserResponse{
		User: *user,
	}, nil
}

func (s *userService) GetUsers(
	ctx context.Context,
	req *model.GetUsersRequest,
) (*model.GetUsersResponse, error) {
	filter := &repository.FindUserByFilter{
		Name:  req.Name,
		Role:  req.Role,
		Page:  req.Page,
		Limit: req.Limit,
	}

	errGroup, errCtx := errgroup.WithContext(ctx)

	var users []entity.User
	errGroup.Go(func() error {
		var err error
		users, err = s.postgresRepo.UserRepo.FindManyByFilter(errCtx, nil, filter)
		if err != nil {
			return err
		}
		return nil
	})

	var count int64
	errGroup.Go(func() error {
		var err error
		count, err = s.postgresRepo.UserRepo.CountByFilter(errCtx, nil, filter)
		if err != nil {
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.GetUsersResponse{
		Users: users,
		Count: count,
	}, nil
}
