package services

import (
	"context"
)

type ArchiveService struct {
	ctx context.Context
}

func NewArchiveService() *ArchiveService {
	return &ArchiveService{}
}

func (s *ArchiveService) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *ArchiveService) CancelOperation() {
}
