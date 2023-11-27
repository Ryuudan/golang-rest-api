package services

import (
	"context"
	"errors"

	"github.com/ryuudan/golang-rest-api/ent/generated"
	"github.com/ryuudan/golang-rest-api/src/internal/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser *generated.User) (*generated.User, error)
	GetUserByID(ctx context.Context, id int) (*generated.User, error)
	GetUserByEmail(ctx context.Context, email string) (*generated.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (user *userService) CreateUser(ctx context.Context, newUser *generated.User) (*generated.User, error) {

	// Check if the email is already taken
	existingUser, err := user.repo.GetByEmail(ctx, newUser.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create the user
	createdUser, err := user.repo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (user *userService) GetUserByID(ctx context.Context, id int) (*generated.User, error) {
	// Additional business logic can be added here before retrieving the user
	return user.repo.GetByID(ctx, id)
}

func (user *userService) GetUserByEmail(ctx context.Context, email string) (*generated.User, error) {
	// Additional business logic can be added here before retrieving the user
	return user.repo.GetByEmail(ctx, email)
}
