package handler

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/pkg/pagination"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
	log              *logger.Logger
}

func NewDashboardHandler(dashboardService *service.DashboardService, log *logger.Logger) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService, log: log}
}

func (h *DashboardHandler) GetSellerDashboard(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	inputPagination := pagination.ParsePagination(r)

	res, err := h.dashboardService.GetSellerDashboard(r.Context(), auth, inputPagination)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *DashboardHandler) GetAdminDashboard(w http.ResponseWriter, r *http.Request) error {

	inputPagination := pagination.ParsePagination(r)

	res, err := h.dashboardService.GetAdminDashboard(r.Context(), inputPagination)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}
