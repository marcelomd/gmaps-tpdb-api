package core

import "github.com/google/uuid"

func NewId() (string, error) {
	id, err := uuid.NewV7()
	return id.String(), err
}
