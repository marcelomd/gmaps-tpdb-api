package user

import (
	"context"
	"fmt"
	"log/slog"

	"fragments/internal/core"
	"fragments/internal/core/interfaces"
	"fragments/internal/core/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(userRepository interfaces.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s UserService) Create(ctx context.Context, nu models.NewUser) (models.User, error) {
	slog.Info("Create", slog.Any("nu", nu))

	if nu.Name == "" {
		return models.User{}, fmt.Errorf("no name")
	}

	if nu.Email == "" {
		return models.User{}, fmt.Errorf("no email")
	}

	if nu.Role == "" {
		return models.User{}, fmt.Errorf("no role")
	}

	if nu.Password == "" {
		return models.User{}, fmt.Errorf("no password")
	}

	id, err := core.NewId()
	if err != nil {
		return models.User{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	input := models.User{
		Id:           id,
		Role:         nu.Role,
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
	}
	output, err := s.userRepository.Create(ctx, input)
	if err != nil {
		return models.User{}, err
	}
	slog.Info("UserService.Create", slog.String("id", id), slog.Any("output", output))
	return output, nil
}

func (s UserService) GetById(ctx context.Context, id string) (models.User, error) {
	return s.userRepository.GetById(ctx, id)
}

func (s UserService) AuthenticateByEmail(ctx context.Context, email string, password string) (models.User, error) {
	u, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}
	err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}
