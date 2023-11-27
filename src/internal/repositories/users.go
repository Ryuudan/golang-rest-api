package repositories

import (
	"context"

	"github.com/ryuudan/golang-rest-api/ent/generated"
	"github.com/ryuudan/golang-rest-api/ent/generated/user"
)

type UserRepository interface {
	Create(ctx context.Context, newUser *generated.User) (*generated.User, error)
	GetByID(ctx context.Context, id int) (*generated.User, error)
	GetByEmail(ctx context.Context, email string) (*generated.User, error)
}

type userRepository struct {
	client *generated.UserClient
}

func NewUserRepository(client *generated.UserClient) UserRepository {
	return &userRepository{client: client}
}

func (repo *userRepository) Create(ctx context.Context, newUser *generated.User) (*generated.User, error) {

	user, err := repo.client.Create().
		SetEmail(newUser.Email).
		SetFirstName(newUser.FirstName).
		SetLastName(newUser.LastName).
		SetPassword(newUser.Password).
		SetNillableMiddleName(newUser.MiddleName).
		SetNillableBirthday(newUser.Birthday).
		SetNillablePhoneNumber(newUser.PhoneNumber).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) GetByID(ctx context.Context, id int) (*generated.User, error) {
	user, err := repo.client.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*generated.User, error) {
	user, err := repo.client.Query().Where(
		user.EmailEQ(email),
	).First(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}
