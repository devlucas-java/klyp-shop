package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/auth"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/utils"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type AuthHandler struct {
	authService *service.AuthService
	log         *logger.Logger
}

func NewAuthHandler(authService *service.AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, log: log}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	res, err := h.authService.Login(&req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req auth.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	res, err := h.authService.Register(&req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	var req auth.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	if err := h.authService.UpdatePassword(&req, auth); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
	return nil
}

func (h *AuthHandler) VerifyPassword(w http.ResponseWriter, r *http.Request) error {
	var req auth.VerifyPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	res, err := h.authService.VerifyPassword(&req, auth)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}
