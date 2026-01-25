package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/pick-cee/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

var (
	InvalidDataRequest = errors.New("Invalid User request data, please recheck.")
	InvalidPassword = errors.New("Incorrect password.")
)

type Service interface {
	CreateUser(ctx context.Context, createUser createUserDto) (repo.User, error)
	LoginUser (ctx context.Context, loginUser loginUserDto) (LoginResponse, error)
}

type svc struct {
	repo *repo.Queries
	db *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{repo: repo, db: db}
}

func (s *svc) CreateUser(ctx context.Context, createUser createUserDto) (repo.User, error) {
	validate := createUser.Validate()

	if !validate {
		return repo.User{}, InvalidDataRequest
	}

	hashedPassword, err := HashPassword(createUser.Password)
	if err != nil {
		return repo.User{}, err
	} 

	createdUser, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		Name: createUser.Name,
		Username: createUser.Username,
		Password: hashedPassword,
	})
	
	if err != nil {
		return repo.User{}, err
	}

	return createdUser, nil
}

func (s *svc) LoginUser(ctx context.Context, loginUser loginUserDto) (LoginResponse, error) {
	user, err := s.repo.FindUserByUsername(ctx, loginUser.Username)

	if err != nil {
		return LoginResponse{}, err
	}

	// compare passwords
	comparedPassword := CheckPassword(loginUser.Password, user.Password)
	if !comparedPassword {
		return LoginResponse{}, InvalidPassword
	}
	
	user.Password = ""
	tokenString, err := CreateToken(int(user.ID), user.Username, user.Name)

	if err != nil {
		return LoginResponse{}, err
	}

	response := LoginResponse{
		User: user,
		Token: tokenString,
	}

	return response, nil
}
