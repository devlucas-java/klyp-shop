package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/product"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/utils"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
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
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	var dto product.CreateProduct
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperrors.BadRequest("invalid request payload", err)
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
		return apperrors.InvalidUUID(err)
	}
	res, err := h.productService.GetProductByID(uuid)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(err)
	}
	var dto product.UpdateProduct
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperrors.BadRequest("invalid request payload", err)
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
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(err)
	}
	if err := h.productService.DeleteProduct(auth, uuid); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}

func (h *ProductHandler) SetTop10(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(err)
	}
	res, err := h.productService.SetTop10(r.Context(), auth, uuid)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}
