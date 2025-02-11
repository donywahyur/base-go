package utils

import "github.com/google/uuid"

type UUIDGenerator interface {
	NewString() string
}

type RealUUIDGenerator struct{}

func (r *RealUUIDGenerator) NewString() string {
	return uuid.NewString()
}
