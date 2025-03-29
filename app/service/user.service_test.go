package service

import (
	"context"
	"testing"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	helper_mocks "sondth-test_soa/mocks/helper"
	repo_mocks "sondth-test_soa/mocks/repository"
	"sondth-test_soa/package/errors"
	"sondth-test_soa/utils"

	"github.com/google/uuid"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	username          = "testuser"
	password          = "testpass"
	hashedPassword, _ = utils.HashPassword(password)
	newPassword       = "newpass"
	fullname          = "Test User"
	role              = entity.ROLE_USER
	adminRole         = entity.ROLE_ADMIN
)

func Test_userService_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.UserRegisterRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.UserRegisterResponse
		wantErr bool
		mock    func(repo *repo_mocks.IUserRepository)
	}

	ctx := context.Background()

	tests := []testCase{
		{
			name: "Register Success",
			args: args{
				ctx: ctx,
				req: &model.UserRegisterRequest{
					Username: username,
					Password: password,
					Fullname: fullname,
				},
			},
			want:    &model.UserRegisterResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IUserRepository) {
				// Mock username check
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Username != nil && *filter.Username == username
				})).Return(nil, gorm.ErrRecordNotFound).Once()

				// Mock user creation
				repo.On("Create", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.Username == username &&
						user.Password == password &&
						user.Fullname == fullname &&
						user.Role == entity.ROLE_USER
				})).Return(nil).Once()
			},
		},
		{
			name: "Register Failed - Username Exists",
			args: args{
				ctx: ctx,
				req: &model.UserRegisterRequest{
					Username: username,
					Password: password,
					Fullname: fullname,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository) {
				// Mock username exists
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Username != nil && *filter.Username == username
				})).Return(&entity.User{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIUserRepository(t)

			// Setup mocks
			tt.mock(repo)

			s := &userService{
				postgresRepo: repository.RepositoryCollections{
					UserRepo: repo,
				},
			}

			got, err := s.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("userService.Register() got nil response, want non-nil")
			}
		})
	}
}

func Test_userService_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.UserLoginRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.UserLoginResponse
		wantErr bool
		mock    func(repo *repo_mocks.IUserRepository, oauthHelper *helper_mocks.IOAuthHelper)
	}

	ctx := context.Background()
	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"

	tests := []testCase{
		{
			name: "Login Success",
			args: args{
				ctx: ctx,
				req: &model.UserLoginRequest{
					Username: username,
					Password: password,
				},
			},
			want: &model.UserLoginResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
			wantErr: false,
			mock: func(repo *repo_mocks.IUserRepository, oauthHelper *helper_mocks.IOAuthHelper) {
				// Mock find user
				user := &entity.User{
					Username: username,
					Password: hashedPassword,
				}
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Username != nil && *filter.Username == username
				})).Return(user, nil).Once()

				// Mock token generation
				oauthHelper.On("GenerateAccessToken", mock.AnythingOfType("entity.User")).Return(accessToken, nil).Once()
				oauthHelper.On("GenerateRefreshToken", mock.AnythingOfType("entity.User")).Return(refreshToken, nil).Once()
			},
		},
		{
			name: "Login Failed - User Not Found",
			args: args{
				ctx: ctx,
				req: &model.UserLoginRequest{
					Username: username,
					Password: password,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository, oauthHelper *helper_mocks.IOAuthHelper) {
				// Mock user not found
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Username != nil && *filter.Username == username
				})).Return(nil, gorm.ErrRecordNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIUserRepository(t)
			oauthHelper := helper_mocks.NewIOAuthHelper(t)

			// Setup mocks
			tt.mock(repo, oauthHelper)

			s := &userService{
				postgresRepo: repository.RepositoryCollections{
					UserRepo: repo,
				},
				helper: helper.HelperCollections{
					OAuthHelper: oauthHelper,
				},
			}

			got, err := s.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("userService.Login() got nil response, want non-nil")
			}
		})
	}
}

func Test_userService_GetUsers(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.GetUsersRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.GetUsersResponse
		wantErr bool
		mock    func(repo *repo_mocks.IUserRepository)
	}

	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID:   userID,
		Role: adminRole,
	})

	users := []entity.User{
		{
			ID:       userID,
			Username: username,
			Fullname: fullname,
			Role:     role,
		},
	}

	tests := []testCase{
		{
			name: "Get Users Success",
			args: args{
				ctx: ctx,
				req: &model.GetUsersRequest{
					Name:  &fullname,
					Role:  &role,
					Page:  &page,
					Limit: &limit,
				},
			},
			want: &model.GetUsersResponse{
				Users: users,
				Count: 1,
			},
			wantErr: false,
			mock: func(repo *repo_mocks.IUserRepository) {
				// Mock find users
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Name != nil && *filter.Name == fullname &&
						filter.Role != nil && *filter.Role == role &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(users, nil).Once()

				// Mock count users
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Name != nil && *filter.Name == fullname &&
						filter.Role != nil && *filter.Role == role &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(int64(1), nil).Once()
			},
		},
		{
			name: "Get Users Failed - Find Error",
			args: args{
				ctx: ctx,
				req: &model.GetUsersRequest{
					Name:  &fullname,
					Role:  &role,
					Page:  &page,
					Limit: &limit,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository) {
				// Mock find users error
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Name != nil && *filter.Name == fullname &&
						filter.Role != nil && *filter.Role == role &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(nil, errors.New(errors.ErrCodeInternalServerError)).Once()

				// Mock count users error
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Name != nil && *filter.Name == fullname &&
						filter.Role != nil && *filter.Role == role &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(int64(0), errors.New(errors.ErrCodeInternalServerError)).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIUserRepository(t)

			// Setup mocks
			tt.mock(repo)

			s := &userService{
				postgresRepo: repository.RepositoryCollections{
					UserRepo: repo,
				},
			}

			got, err := s.GetUsers(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("userService.GetUsers() got nil response, want non-nil")
			}
		})
	}
}

func Test_userService_ForgetPassword(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.ForgetPasswordRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.ForgetPasswordResponse
		wantErr bool
		mock    func(repo *repo_mocks.IUserRepository)
	}

	ctx := context.Background()

	tests := []testCase{
		{
			name: "Forget Password Success",
			args: args{
				ctx: ctx,
				req: &model.ForgetPasswordRequest{
					Username:    username,
					NewPassword: newPassword,
				},
			},
			want:    &model.ForgetPasswordResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IUserRepository) {
				// Mock find user
				user := &entity.User{
					Username: username,
					Password: hashedPassword,
				}
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Username != nil && *filter.Username == username
				})).Return(user, nil).Once()

				// Mock update user
				repo.On("Update", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Username == username && u.Password == newPassword
				})).Return(nil).Once()
			},
		},
		{
			name: "Forget Password Failed - User Not Found",
			args: args{
				ctx: ctx,
				req: &model.ForgetPasswordRequest{
					Username:    username,
					NewPassword: newPassword,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository) {
				// Mock user not found
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.Username != nil && *filter.Username == username
				})).Return(nil, gorm.ErrRecordNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIUserRepository(t)

			// Setup mocks
			tt.mock(repo)

			s := &userService{
				postgresRepo: repository.RepositoryCollections{
					UserRepo: repo,
				},
			}

			got, err := s.ForgetPassword(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.ForgetPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("userService.ForgetPassword() got nil response, want non-nil")
			}
		})
	}
}

func Test_userService_ChangePassword(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.ChangeUserPasswordRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.ChangeUserPasswordResponse
		wantErr bool
		mock    func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper)
	}

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID: userID,
	})

	tests := []testCase{
		{
			name: "Change Password Success",
			args: args{
				ctx: ctx,
				req: &model.ChangeUserPasswordRequest{
					ID:          userID,
					OldPassword: password,
					NewPassword: newPassword,
				},
			},
			want:    &model.ChangeUserPasswordResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock permission check
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(true).Once()

				// Mock find user
				user := &entity.User{
					ID:       userID,
					Password: hashedPassword,
				}
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.ID != nil && *filter.ID == userID
				})).Return(user, nil).Once()

				// Mock update user
				repo.On("Update", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.ID == userID && u.Password == newPassword
				})).Return(nil).Once()
			},
		},
		{
			name: "Change Password Failed - Invalid Permission",
			args: args{
				ctx: ctx,
				req: &model.ChangeUserPasswordRequest{
					ID:          userID,
					OldPassword: password,
					NewPassword: newPassword,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock invalid permission
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(false).Once()
			},
		},
		{
			name: "Change Password Failed - User Not Found",
			args: args{
				ctx: ctx,
				req: &model.ChangeUserPasswordRequest{
					ID:          userID,
					OldPassword: password,
					NewPassword: newPassword,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock permission check
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(true).Once()

				// Mock user not found
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.ID != nil && *filter.ID == userID
				})).Return(nil, gorm.ErrRecordNotFound).Once()
			},
		},
		{
			name: "Change Password Failed - Incorrect Old Password",
			args: args{
				ctx: ctx,
				req: &model.ChangeUserPasswordRequest{
					ID:          userID,
					OldPassword: "wrongpass",
					NewPassword: newPassword,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock permission check
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(true).Once()

				// Mock find user
				user := &entity.User{
					ID:       userID,
					Password: hashedPassword,
				}
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.ID != nil && *filter.ID == userID
				})).Return(user, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIUserRepository(t)
			userHelper := helper_mocks.NewIUserHelper(t)

			// Setup mocks
			tt.mock(repo, userHelper)

			s := &userService{
				postgresRepo: repository.RepositoryCollections{
					UserRepo: repo,
				},
				helper: helper.HelperCollections{
					UserHelper: userHelper,
				},
			}

			got, err := s.ChangePassword(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("userService.ChangePassword() got nil response, want non-nil")
			}
		})
	}
}

func Test_userService_UpdateUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.UpdateUserRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.UpdateUserResponse
		wantErr bool
		mock    func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper)
	}

	userID := uuid.New()
	newUsername := "newuser"
	newFullname := "New User"
	adminRole := entity.ROLE_ADMIN

	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID:   userID,
		Role: adminRole,
	})

	tests := []testCase{
		{
			name: "Update User Success",
			args: args{
				ctx: ctx,
				req: &model.UpdateUserRequest{
					ID:       userID,
					Username: &newUsername,
					Fullname: &newFullname,
					Role:     &role,
				},
			},
			want: &model.UpdateUserResponse{
				User: entity.User{
					ID:       userID,
					Username: newUsername,
					Fullname: newFullname,
					Role:     role,
				},
			},
			wantErr: false,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock permission check
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(true).Once()

				// Mock find user
				user := &entity.User{
					ID:       userID,
					Username: username,
					Fullname: fullname,
					Role:     role,
				}
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.ID != nil && *filter.ID == userID
				})).Return(user, nil).Once()

				// Mock role validation
				userHelper.On("IsValidRole", role).Return(true).Once()

				// Mock update user
				repo.On("Update", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.ID == userID &&
						u.Username == newUsername &&
						u.Fullname == newFullname &&
						u.Role == role
				})).Return(nil).Once()
			},
		},
		{
			name: "Update User Failed - Invalid Permission",
			args: args{
				ctx: ctx,
				req: &model.UpdateUserRequest{
					ID:       userID,
					Username: &newUsername,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock invalid permission
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(false).Once()
			},
		},
		{
			name: "Update User Failed - No User Context",
			args: args{
				ctx: context.Background(),
				req: &model.UpdateUserRequest{
					ID:       userID,
					Username: &newUsername,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock no user context
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(false).Once()
			},
		},
		{
			name: "Update User Failed - User Not Found",
			args: args{
				ctx: ctx,
				req: &model.UpdateUserRequest{
					ID:       userID,
					Username: &newUsername,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock permission check
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(true).Once()

				// Mock user not found
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.ID != nil && *filter.ID == userID
				})).Return(nil, gorm.ErrRecordNotFound).Once()
			},
		},
		{
			name: "Update User Failed - Invalid Role",
			args: args{
				ctx: ctx,
				req: &model.UpdateUserRequest{
					ID:   userID,
					Role: &role,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IUserRepository, userHelper *helper_mocks.IUserHelper) {
				// Mock permission check
				userHelper.On("IsValidPermission", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), userID).Return(true).Once()

				// Mock find user
				user := &entity.User{
					ID:       userID,
					Username: username,
					Fullname: fullname,
					Role:     role,
				}
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindUserByFilter) bool {
					return filter.ID != nil && *filter.ID == userID
				})).Return(user, nil).Once()

				// Mock invalid role
				userHelper.On("IsValidRole", role).Return(false).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIUserRepository(t)
			userHelper := helper_mocks.NewIUserHelper(t)

			// Setup mocks
			tt.mock(repo, userHelper)

			s := &userService{
				postgresRepo: repository.RepositoryCollections{
					UserRepo: repo,
				},
				helper: helper.HelperCollections{
					UserHelper: userHelper,
				},
			}

			got, err := s.UpdateUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("userService.UpdateUser() got nil response, want non-nil")
			}
		})
	}
}

func TestNewUserService(t *testing.T) {
	type args struct {
		postgresRepo repository.RepositoryCollections
		helper       helper.HelperCollections
	}
	tests := []struct {
		name string
		args args
		want IUserService
	}{
		{
			name: "NewUserService",
			args: args{
				postgresRepo: repository.RepositoryCollections{
					UserRepo: repo_mocks.NewIUserRepository(t),
				},
				helper: helper.HelperCollections{
					UserHelper:  helper_mocks.NewIUserHelper(t),
					OAuthHelper: helper_mocks.NewIOAuthHelper(t),
				},
			},
			want: &userService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserService(tt.args.postgresRepo, tt.args.helper)
			if got == nil {
				t.Error("NewUserService() got nil service")
			}
		})
	}
}
