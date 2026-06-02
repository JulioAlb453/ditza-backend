package domain

import "fmt"

type Status string

const (
	StatusPending  Status = "pending"
	StatusAccepted Status = "accepted"
	StatusRejected Status = "rejected"
)

func ParseStatus(raw string) (Status, error) {
	status := Status(raw)
	switch status {
	case StatusPending, StatusAccepted, StatusRejected:
		return status, nil
	default:
		return "", fmt.Errorf("estado de amistad inválido: %s", raw)
	}
}

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	_, err := ParseStatus(string(s))
	return err == nil
}

func (s Status) CanAccept() bool {
	return s == StatusPending
}

func (s Status) CanReject() bool {
	return s == StatusPending
}
