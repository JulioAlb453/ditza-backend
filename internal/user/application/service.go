package application

import (
	"context"
	"strings"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
	"ditza/internal/shared/infrastructure/password"
	userdomain "ditza/internal/user/domain"
)

type Service struct {
	userRepository userdomain.Repository
}

type RegisterCommand struct {
	Alias    string
	Email    string
	Password string
}

type RegisterResult struct {
	UserID valueobject.UserID
	Alias  string
	Email  string
}

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	UserID valueobject.UserID
	Alias  string
	Email  string
}

func NewService(userRepository userdomain.Repository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Register(ctx context.Context, command RegisterCommand) (result *RegisterResult, err error) {
	tracker := monitoring.StartService(logger.ModelUser, "register", map[string]any{"email": command.Email})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"user_id": result.UserID})
	}()

	email := strings.TrimSpace(strings.ToLower(command.Email))
	exists, err := s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerror.New("USER_ALREADY_EXISTS", "el correo ya está registrado", domainerror.ErrInvalidInput)
	}

	hashedPassword, err := password.Hash(command.Password)
	if err != nil {
		return nil, err
	}

	userEntity, err := userdomain.New(
		valueobject.NewUserID(),
		command.Alias,
		email,
		hashedPassword,
	)
	if err != nil {
		return nil, err
	}

	if err := s.userRepository.Create(ctx, userEntity); err != nil {
		return nil, err
	}

	return &RegisterResult{
		UserID: userEntity.ID,
		Alias:  userEntity.Alias,
		Email:  userEntity.Email,
	}, nil
}

func (s *Service) Login(ctx context.Context, command LoginCommand) (result *LoginResult, err error) {
	tracker := monitoring.StartService(logger.ModelUser, "login", map[string]any{"email": command.Email})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"user_id": result.UserID})
	}()

	email := strings.TrimSpace(strings.ToLower(command.Email))
	userEntity, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if userEntity == nil || !password.Verify(userEntity.Password, command.Password) {
		return nil, domainerror.New("INVALID_CREDENTIALS", "correo o contraseña incorrectos", domainerror.ErrUnauthorized)
	}

	return &LoginResult{
		UserID: userEntity.ID,
		Alias:  userEntity.Alias,
		Email:  userEntity.Email,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, userID valueobject.UserID) (userEntity *userdomain.User, err error) {
	tracker := monitoring.StartService(logger.ModelUser, "get_by_id", map[string]any{"user_id": userID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(nil)
	}()

	userEntity, err = s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userEntity == nil {
		return nil, domainerror.New("USER_NOT_FOUND", "usuario no encontrado", domainerror.ErrNotFound)
	}
	return userEntity, nil
}

const defaultSearchLimit = 20

type SearchUsersResult struct {
	UserID valueobject.UserID
	Alias  string
}

func (s *Service) SearchByAlias(ctx context.Context, requesterID valueobject.UserID, aliasQuery string) (results []SearchUsersResult, err error) {
	tracker := monitoring.StartService(logger.ModelUser, "search_by_alias", map[string]any{"alias": aliasQuery})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"count": len(results)})
	}()

	aliasQuery = strings.TrimSpace(aliasQuery)
	if len(aliasQuery) < 2 {
		return nil, domainerror.New("INVALID_SEARCH_QUERY", "el alias debe tener al menos 2 caracteres", domainerror.ErrInvalidInput)
	}

	users, err := s.userRepository.SearchByAlias(ctx, aliasQuery, requesterID, defaultSearchLimit)
	if err != nil {
		return nil, err
	}

	results = make([]SearchUsersResult, 0, len(users))
	for _, userEntity := range users {
		results = append(results, SearchUsersResult{
			UserID: userEntity.ID,
			Alias:  userEntity.Alias,
		})
	}
	return results, nil
}
