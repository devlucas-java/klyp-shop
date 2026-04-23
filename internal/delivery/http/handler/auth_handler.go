package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dauth"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type AuthHandler struct {
	authService *service.AuthService
	log         *logger.Logger
}

func NewAuthHandler(authService *service.AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {

	var req dauth.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("invalid request payload", err)
	}

	res, err := h.authService.Login(&req)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req dauth.RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warnf("Invalid JSON payload for register: %v", err)
		return errors.ErrBadRequest("invalid request payload", err)
	}

	res, err := h.authService.Register(&req)
	if err != nil {
		h.log.Errorf("Registration failed for duser %s: %v", req.Username, err)
		return err
	}

	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) error {

	var req dauth.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warnf("Invalid JSON payload for change password: %v", err)
		return errors.ErrBadRequest("invalid request payload", err)
	}

	user := r.Context().Value(middleware.AuthKey).(*entity.User)

	err := h.authService.UpdatePassword(&req, user)
	if err != nil {
		h.log.Errorf("Password update failed: %v", err)
		return err
	}

	response.ResponseEntity(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
	return nil
}

func (h *AuthHandler) VerifyPassword(w http.ResponseWriter, r *http.Request) error {
	var req dauth.VerifyPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return errors.ErrBadRequest("invalid request payload", err)
	}

	if req.Password == "" {
		return errors.ErrBadRequest("password is required", nil)
	}

	user := r.Context().Value(middleware.AuthKey).(*entity.User)

	res, err := h.authService.VerifyPassword(&req, user)
	if err != nil {
		h.log.Warnf("Password verification failed: %v", err)
		return err
	}

	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}
