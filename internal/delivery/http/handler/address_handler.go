package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type AddressHandler struct {
	addressService *service.AddressService
	log            *logger.Logger
}

func NewAddressHandler(addressService *service.AddressService, log *logger.Logger) *AddressHandler {
	return &AddressHandler{addressService: addressService, log: log}
}

func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) error {

	var req daddress.CreateAddressRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("Payload invalid", err)
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

	var req daddress.UpdateAddressRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("Payload invalid", err)
	}

	idStr := chi.URLParam(r, "id")

	id, err := id.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Id invalid", err)
	}

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	res, err := h.addressService.UpdateAddress(auth, &req, id)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) error {

	idStr := chi.URLParam(r, "id")

	id, err := id.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Id invalid", err)
	}

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	err = h.addressService.DeleteAddress(auth, id)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusOK, nil)
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
