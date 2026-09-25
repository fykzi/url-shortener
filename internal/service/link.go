package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/fykzi/url-shortener/internal/model"
	"github.com/fykzi/url-shortener/internal/repository"
	"github.com/fykzi/url-shortener/internal/repository/sqlite"
	"github.com/fykzi/url-shortener/internal/shortener"
)

var (
	ErrInvalidUrl = errors.New("invalid url")

	ErrCodeCollision = errors.New("code collision")

	ErrNotFound = errors.New("link not found")
)

const (
	maxCodeAttempts = 5
	codeLength      = 6
)

type LinkService struct {
	repo    repository.LinkRepository
	BaseURL string
	log     *slog.Logger
}

func NewLinkService(repo repository.LinkRepository, url string, log *slog.Logger) LinkService {
	return LinkService{
		repo:    repo,
		BaseURL: url,
		log:     log,
	}
}

func (s *LinkService) Create(ctx context.Context, rawURL string) (*model.Link, error) {
	if err := validateURL(rawURL); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidUrl, err)
	}

	for range maxCodeAttempts {
		code, err := shortener.GenerateCode(codeLength)
		if err != nil {
			return nil, fmt.Errorf("generate code: %w", err)
		}

		now := time.Now().UTC()
		link := &model.Link{
			ShortCode:   code,
			OriginalURL: rawURL,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		err = s.repo.Create(ctx, link)
		if err == nil {
			return link, nil
		}

		if isUniqueViolantion(err) {
			continue
		}

		return nil, fmt.Errorf("create link: %w", err)
	}

	return nil, ErrCodeCollision
}

func (s *LinkService) Redirect(ctx context.Context, code string) (string, error) {
	originalURL, err := s.repo.Redirect(ctx, code)
	if err != nil {
		if errors.Is(err, sqlite.ErrNotFound) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("redirect: get link: %w", err)
	}

	return originalURL, nil
}

func (s *LinkService) GetByCode(ctx context.Context, code string) (*model.Link, error) {
	link, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sqlite.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get link by code: %w", err)
	}

	return link, nil
}

func (s *LinkService) Delete(ctx context.Context, code string) error {
	err := s.repo.Delete(ctx, code)
	if err != nil {
		if errors.Is(err, sqlite.ErrNotFound) {
			return ErrNotFound
		}
	}

	return err
}

func (s *LinkService) ShortURL(code string) string {
	return s.BaseURL + code
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return errors.New("empty url")
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return errors.New("malformed url")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme must be http or https")
	}

	if u.Host == "" {
		return errors.New("missing host")
	}

	return nil
}

func isUniqueViolantion(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
