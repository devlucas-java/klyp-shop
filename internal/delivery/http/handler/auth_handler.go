package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/request/auth_request"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth_request.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warnf("Invalid JSON payload for login: %v", err)
		h.respondWithError(w, errors.New("INVALID_PAYLOAD", "invalid request payload", http.StatusBadRequest, err))
		return
	}

	res, err := h.authService.Login(&req)
	if err != nil {
		h.log.Warnf("Login failed for user %s: %v", req.Login, err)
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth_request.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warnf("Invalid JSON payload for register: %v", err)
		h.respondWithError(w, errors.New("INVALID_PAYLOAD", "invalid request payload", http.StatusBadRequest, err))
		return
	}

	res, err := h.authService.Register(&req)
	if err != nil {
		h.log.Errorf("Registration failed for user %s: %v", req.Username, err)
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, res)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {

	var req auth_request.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warnf("Invalid JSON payload for change password: %v", err)
		h.respondWithError(w, errors.New("INVALID_PAYLOAD", "invalid request payload", http.StatusBadRequest, err))
		return
	}

	user := r.Context().Value(middleware.AuthKey).(*entity.User)

	err := h.authService.UpdatePassword(&req, user)
	if err != nil {
		h.log.Errorf("Password update failed: %v", err)
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
}

func (h *AuthHandler) VerifyPassword(w http.ResponseWriter, r *http.Request) {
	var req auth_request.VerifyPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.respondWithError(w, errors.New("INVALID_BODY", "invalid request body", http.StatusBadRequest, nil))
		return
	}

	if req.Password == "" {
		h.respondWithError(w, errors.New("MISSING_PASSWORD", "password is required", http.StatusBadRequest, nil))
		return
	}

	user := r.Context().Value(middleware.AuthKey).(*entity.User)

	res, err := h.authService.VerifyPassword(&req, user)
	if err != nil {
		h.log.Warnf("Password verification failed: %v", err)
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) respondWithError(w http.ResponseWriter, err error) {
	var appErr *errors.AppError
	switch e := err.(type) {
	case *errors.AppError:
		appErr = e
	default:
		appErr = errors.New("INTERNAL_ERROR", "an unexpected error occurred", http.StatusInternalServerError, err)
	}

	h.respondWithJSON(w, appErr.Status, appErr)
}

func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		h.log.Errorf("Error marshaling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code":"INTERNAL_ERROR","message":"error marshaling response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
