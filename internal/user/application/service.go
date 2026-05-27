package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	userdomain "ditza/internal/user/domain"
)

type Service struct {
	userRepository userdomain.Repository
}

type RegisterCommand struct {
	Name         string
	Email        string
	PasswordHash string
	Timezone     string
}

type RegisterResult struct {
	UserID     valueobject.UserID
	Name       string
	Email      string
	Timezone   string
	FriendCode string
}

func NewService(userRepository userdomain.Repository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Register(ctx context.Context, command RegisterCommand) (*RegisterResult, error) {
	email := strings.TrimSpace(strings.ToLower(command.Email))
	exists, err := s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerror.New("USER_ALREADY_EXISTS", "el correo ya está registrado", domainerror.ErrInvalidInput)
	}

	friendCode, err := generateFriendCode()
	if err != nil {
		return nil, err
	}

	userEntity, err := userdomain.New(
		command.Name,
		email,
		command.PasswordHash,
		command.Timezone,
		friendCode,
	)
	if err != nil {
		return nil, err
	}

	if err := s.userRepository.Create(ctx, userEntity); err != nil {
		return nil, err
	}

	return &RegisterResult{
		UserID:     userEntity.ID,
		Name:       userEntity.Name,
		Email:      userEntity.Email,
		Timezone:   userEntity.Timezone,
		FriendCode: userEntity.FriendCode,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, userID valueobject.UserID) (*userdomain.User, error) {
	userEntity, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userEntity == nil {
		return nil, domainerror.New("USER_NOT_FOUND", "usuario no encontrado", domainerror.ErrNotFound)
	}
	return userEntity, nil
}

func generateFriendCode() (string, error) {
	bytes := make([]byte, valueobject.FriendCodeLength/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", domainerror.New("FRIEND_CODE_GENERATION_FAILED", "no se pudo generar el código de amigo", err)
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}
