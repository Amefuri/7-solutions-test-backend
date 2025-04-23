package user

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, name, email, password string) (*User, error) {
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("error getting user by email: " + err.Error())
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &User{Name: name, Email: email, Password: string(hashed), CreatedAt: time.Now()}
	return user, s.repo.Create(ctx, user)
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	user, _ := s.repo.GetByEmail(ctx, email)
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, user *User) error {
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
