package valueobject

import "github.com/google/uuid"

type UserID string

func NewUserID() UserID {
	return UserID(uuid.NewString())
}

func ParseUserID(raw string) (UserID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return "", err
	}
	return UserID(parsed.String()), nil
}

func (id UserID) String() string {
	return string(id)
}

func (id UserID) IsZero() bool {
	return id == ""
}
