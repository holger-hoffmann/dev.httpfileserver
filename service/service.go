package service

import (
	"context"
)

type Service struct {
	rootDir string
}

func New(ctx context.Context, rootDir string) (*Service, error) {
	return &Service{
		rootDir: rootDir,
	}, nil
}
