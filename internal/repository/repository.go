package repository

import (
	"context"

	"github.com/fykzi/url-shortener/internal/model"
)

type LinkRepository interface {
	Create(ctx context.Context, link *model.Link) error
	GetByCode(ctx context.Context, code string) (*model.Link, error)
	Delete(ctx context.Context, code string) error
	List(ctx context.Context) ([]*model.Link, error)
	Redirect(ctx context.Context, code string) (string, error)
}
