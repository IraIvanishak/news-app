package services

import (
	"context"

	"github.com/IraIvanishak/news-app/internal/core/domain"
	"github.com/IraIvanishak/news-app/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	repo ports.PostRepository
}

var _ ports.PostService = (*PostService)(nil)

func NewPostService(repo ports.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) Store(ctx context.Context, m *domain.Post) error {
	return s.repo.Store(ctx, m)
}

func (s *PostService) Find(ctx context.Context, id primitive.ObjectID) (*domain.Post, error) {
	return s.repo.Find(ctx, id)
}

func (s *PostService) List(ctx context.Context, filter domain.ListFilter) ([]*domain.Post, error) {
	return s.repo.List(ctx, filter)
}

func (s *PostService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *PostService) Update(ctx context.Context, m *domain.Post) error {
	return s.repo.Update(ctx, m)
}

func (s *PostService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
