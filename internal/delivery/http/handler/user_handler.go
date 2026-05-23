package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/user"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/utils"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	userService *service.UserService
	log         *logger.Logger
}

func NewUserHandler(userService *service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{userService: userService, log: log}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	res, err := h.userService.GetMe(auth)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	var dto user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperrors.BadRequest("invalid request payload", err)
	}
	if err := dto.Validate(); err != nil {
		return err
	}
	res, err := h.userService.UpdateMe(auth, &dto)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	if err := h.userService.DeleteMe(auth); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}

func (h *UserHandler) PromoteUser(w http.ResponseWriter, r *http.Request) error {
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(err)
	}
	if err := h.userService.PromoteToAdmin(uuid); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}

func (h *UserHandler) DemoteUser(w http.ResponseWriter, r *http.Request) error {
	uuid, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(err)
	}
	if err := h.userService.DemoteToUser(uuid); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
