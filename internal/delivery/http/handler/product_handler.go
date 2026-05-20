package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dproduct"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	productService *service.ProductService
	log            *logger.Logger
}

func NewProductHandler(productService *service.ProductService, log *logger.Logger) *ProductHandler {
	return &ProductHandler{productService: productService, log: log}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	var dto dproduct.CreateProduct
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return errors.ErrInvalidPayload(err)
	}
	if err := dto.Validate(); err != nil {
		return err
	}
	res, err := h.productService.CreateProduct(auth, &dto)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) error {
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	res, err := h.productService.GetProductByID(uuid)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	var dto dproduct.UpdateProduct
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return errors.ErrInvalidPayload(err)
	}
	if err := dto.Validate(); err != nil {
		return err
	}
	res, err := h.productService.UpdateProduct(auth, &dto, uuid)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	if err := h.productService.DeleteProduct(auth, uuid); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
