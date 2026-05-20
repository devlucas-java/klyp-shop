package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dseller"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type SellerHandler struct {
	sellerService *service.SellerService
	log           *logger.Logger
}

func NewSellerHandler(sellerService *service.SellerService, log *logger.Logger) *SellerHandler {
	return &SellerHandler{sellerService: sellerService, log: log}
}

func (h *SellerHandler) CreateSeller(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	var dto dseller.CreateSeller
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return errors.ErrInvalidPayload(err)
	}
	if err := dto.Validate(); err != nil {
		return err
	}
	res, err := h.sellerService.CreateSeller(auth, &dto)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *SellerHandler) GetSellerByID(w http.ResponseWriter, r *http.Request) error {
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	res, err := h.sellerService.GetSellerByID(uuid)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *SellerHandler) UpdateSeller(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	var dto dseller.UpdateSeller
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return errors.ErrInvalidPayload(err)
	}
	if err := dto.Validate(); err != nil {
		return err
	}
	res, err := h.sellerService.UpdateSeller(auth, &dto)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *SellerHandler) DeleteSeller(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	if err := h.sellerService.DeleteSeller(auth); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
