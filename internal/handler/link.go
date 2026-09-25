package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fykzi/url-shortener/internal/model"
	"github.com/fykzi/url-shortener/internal/service"
	"github.com/fykzi/url-shortener/internal/shortener"
)

const maxRequestBody = 1 << 20

type LinkHandler struct {
	svc service.LinkService
	log *slog.Logger
}

type LinkService interface {
	Create(ctx context.Context, rawURL string) (*model.Link, error)
	GetByCode(ctx context.Context, code string) (*model.Link, error)
	Delete(ctx context.Context, code string) error
	List(ctx context.Context) ([]*model.Link, error)
	Redirect(ctx context.Context, code string) (string, error)
	ShortURL(code string) string
}

func NewLinkHandler(svc service.LinkService, log *slog.Logger) *LinkHandler {
	return &LinkHandler{
		svc: svc,
		log: log,
	}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)

	var req createLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	link, err := h.svc.Create(r.Context(), req.URL)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUrl):
			h.respondError(w, http.StatusBadRequest, "invalid url")
		case errors.Is(err, service.ErrCodeCollision):
			h.log.Error("code collision after retries", "error", err)
			h.respondError(w, http.StatusInternalServerError, "internal error")
		default:
			h.log.Error("create link failed", "error", err)
			h.respondError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	h.log.Info("link created", "code", link.ShortCode)
	h.respondJSON(w, http.StatusCreated, toLinkResponse(link, h.svc.ShortURL(link.ShortCode)))
}

func (h *LinkHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if err := shortener.ValidateCode(code); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	link, err := h.svc.GetByCode(r.Context(), code)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, shortener.ErrInvalidCode) {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			h.respondError(w, http.StatusNotFound, "link not found")
			return
		}
		h.log.Error("get link failed", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.respondJSON(w, http.StatusOK, toLinkResponse(link, h.svc.ShortURL(link.ShortCode)))
}

func (h *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if err := shortener.ValidateCode(code); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(r.Context(), code); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.respondError(w, http.StatusNotFound, "link not found")
			return
		}

		h.log.Error("delete link failed", "code", code, "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.log.Info("delete link", "code", code)
	w.WriteHeader(http.StatusNoContent)
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if err := shortener.ValidateCode(code); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	originalURL, err := h.svc.Redirect(r.Context(), code)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, service.ErrNotFound) {
			h.respondError(w, http.StatusNotFound, "link not found")
			return
		}

		h.log.Error("redirect failed", "code", code, "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
