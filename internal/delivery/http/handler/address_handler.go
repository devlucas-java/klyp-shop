package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	addressDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/address"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

const addressHandler = "address_handler"

type AddressHandler struct {
	addressService *service.AddressService
	log            *logger.Logger
}

func NewAddressHandler(addressService *service.AddressService, log *logger.Logger) *AddressHandler {
	return &AddressHandler{addressService: addressService, log: log}
}

func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) error {
	var req addressDTO.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest(addressHandler+".create_address: invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	res, err := h.addressService.CreateAddress(auth, &req)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) error {
	var req addressDTO.UpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest(addressHandler+".update_address: invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}

	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(addressHandler+".update_address: invalid uuid", err)
	}

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	res, err := h.addressService.UpdateAddress(auth, &req, uuid)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) error {
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(addressHandler+".delete_address: invalid uuid", err)
	}

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	if err := h.addressService.DeleteAddress(auth, uuid); err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusNoContent, nil)
	return nil
}

func (h *AddressHandler) GetAddresses(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	res, err := h.addressService.GetAddresses(auth)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}
