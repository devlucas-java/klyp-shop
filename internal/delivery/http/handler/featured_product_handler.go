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

const featuredProductHandlerTrace = "featured_product_handler.FeaturedProductHandler"

type FeaturedProductHandler struct {
	featuredService *service.FeaturedProductService
	log             *logger.Logger
}

func NewFeaturedProductHandler(featuredService *service.FeaturedProductService, log *logger.Logger) *FeaturedProductHandler {
	return &FeaturedProductHandler{featuredService: featuredService, log: log}
}

func (h *FeaturedProductHandler) AddFeatured(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}

	var req product.AddFeaturedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest(featuredProductHandlerTrace+".add_featured: invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	res, err := h.featuredService.AddFeatured(auth, &req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *FeaturedProductHandler) RemoveFeatured(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}

	productID, err := id.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		return apperrors.InvalidUUID(featuredProductHandlerTrace+".remove_featured: invalid product id", err)
	}
	if err := h.featuredService.RemoveFeatured(auth, productID); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}

func (h *FeaturedProductHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}

	productID, err := id.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		return apperrors.InvalidUUID(featuredProductHandlerTrace+".update_position: invalid product id", err)
	}
	var req product.UpdateFeaturedPositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest(featuredProductHandlerTrace+".update_position: invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	if err := h.featuredService.UpdatePosition(auth, productID, &req); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}

func (h *FeaturedProductHandler) GetMyFeatured(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}

	res, err := h.featuredService.GetMyFeatured(auth)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *FeaturedProductHandler) GetAllFeatured(w http.ResponseWriter, r *http.Request) error {
	res, err := h.featuredService.GetAllFeatured()
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *FeaturedProductHandler) GetFeaturedBySeller(w http.ResponseWriter, r *http.Request) error {
	sellerID, err := id.Parse(chi.URLParam(r, "sellerID"))
	if err != nil {
		return apperrors.InvalidUUID(featuredProductHandlerTrace+".get_featured_by_seller: invalid seller id", err)
	}
	res, err := h.featuredService.GetFeaturedBySeller(sellerID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}
