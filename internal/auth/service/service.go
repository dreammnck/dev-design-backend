package service

import (
	"backend/internal/auth"
	"backend/internal/auth/repository"
	jwtPkg "backend/pkg/jwt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user account is inactive")
)

type AuthService interface {
	Login(req auth.LoginReq) (*auth.LoginRes, error)
	GetUser(id string) (*auth.UserInfo, error)
	GetAllUsers() ([]auth.UserInfo, error)
	AdminUpdateUser(id string, req auth.UpdateUserReq) error
	Register(req auth.RegisterReq) (*auth.UserInfo, error)
	UpdateRole(id string, req auth.UpdateRoleReq) error
	DeleteUser(id string) error
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) GetUser(id string) (*auth.UserInfo, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	return &auth.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *authService) GetAllUsers() ([]auth.UserInfo, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var userInfos []auth.UserInfo
	for _, u := range users {
		status := "active"
		if !u.IsActive {
			status = "suspended"
		}
		userInfos = append(userInfos, auth.UserInfo{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
			Status:   status,
			JoinedAt: u.CreatedAt,
		})
	}

	return userInfos, nil
}

func (s *authService) AdminUpdateUser(id string, req auth.UpdateUserReq) error {
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Role != "" {
		if req.Role != auth.RoleAdmin && req.Role != auth.RoleOrganization && req.Role != auth.RoleCustomer {
			return errors.New("invalid role")
		}
		updates["role"] = req.Role
	}
	if req.Status != "" {
		if req.Status == "active" {
			updates["is_active"] = true
		} else if req.Status == "suspended" {
			updates["is_active"] = false
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return s.repo.Update(id, updates)
}

func (s *authService) Login(req auth.LoginReq) (*auth.LoginRes, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := jwtPkg.GenerateToken(auth.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		return nil, err
	}

	return &auth.LoginRes{
		Token: token,
		User: auth.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (s *authService) Register(req auth.RegisterReq) (*auth.UserInfo, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := auth.RoleCustomer
	if req.Role != "" {
		if req.Role != auth.RoleAdmin && req.Role != auth.RoleOrganization && req.Role != auth.RoleCustomer {
			return nil, errors.New("invalid role")
		}
		role = req.Role
	}

	user := &auth.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsActive:     true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return &auth.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *authService) UpdateRole(id string, req auth.UpdateRoleReq) error {
	if req.Role != auth.RoleAdmin && req.Role != auth.RoleOrganization && req.Role != auth.RoleCustomer {
		return errors.New("invalid role")
	}
	return s.repo.UpdateRole(id, req.Role)
}

func (s *authService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
