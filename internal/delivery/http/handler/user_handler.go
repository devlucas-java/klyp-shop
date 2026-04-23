package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/duser"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	userService *service.UserService
	log         *logger.Logger
}

func NewUserHandler(userService *service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) error {

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	userDTO, err := h.userService.GetMe(auth)
	if err != nil {
		h.log.Errorf("Failed to get duser by ID %s: %v", auth.ID, err)
		return err
	}
	response.ResponseEntity(w, http.StatusOK, userDTO)
	return nil
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) error {

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	var dto duser.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.log.Errorf("Failed to read request: %v", err)
		return errors.ErrBadRequest("invalid request payload", err)
	}

	userDTO, err := h.userService.UpdateMe(auth, &dto)
	if err != nil {
		h.log.Errorf("Failed to update duser profile %s: %v", auth.ID, err)
		return err
	}
	response.ResponseEntity(w, http.StatusOK, userDTO)
	return nil
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) error {

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	err := h.userService.DeleteMe(auth)
	if err != nil {
		h.log.Errorf("Failed to delete duser by ID %s: %v", auth.ID, err)
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}

func (h *UserHandler) PromoteUser(w http.ResponseWriter, r *http.Request) error {

	idString := chi.URLParam(r, "id")
	uuid, err := id.Parse(idString)

	if err != nil {
		h.log.Errorf("Failed to parse duser ID: %s: %v", idString, err)
		return errors.ErrInvalidUUID(err)
	}

	err = h.userService.PromoteToAdmin(uuid)

	if err != nil {
		h.log.Errorf("Failed to promote duser: %s to admin: %v", idString, err)
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}

func (h *UserHandler) DemoteUser(w http.ResponseWriter, r *http.Request) error {

	idString := chi.URLParam(r, "id")
	uuid, err := id.Parse(idString)

	if err != nil {
		h.log.Errorf("Failed to parse duser ID: %s: %v", idString, err)
		return errors.ErrInvalidUUID(err)
	}

	err = h.userService.DemoteToUser(uuid)
	if err != nil {
		h.log.Errorf("Failed to demote duser: %s to admin: %v", idString, err)
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
