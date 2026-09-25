package handler

import (
	"time"

	"github.com/fykzi/url-shortener/internal/model"
)

func toLinkResponse(link *model.Link, shortURL string) linkResponse {
	return linkResponse{
		ShortCode:   link.ShortCode,
		ShortURL:    shortURL,
		OriginalURL: link.OriginalURL,
		Clicks:      link.Clicks,
		CreatedAt:   link.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   link.UpdatedAt.Format(time.RFC3339),
	}
}
