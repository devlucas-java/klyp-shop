package handler

import (
	"net/http"
	"strconv"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
	log              *logger.Logger
}

func NewDashboardHandler(dashboardService *service.DashboardService, log *logger.Logger) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService, log: log}
}

// GET /dashboard/seller?page=1&size=20&status=pending
func (h *DashboardHandler) GetSellerDashboard(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	page, size := parsePagination(r)
	status := r.URL.Query().Get("status")

	res, err := h.dashboardService.GetSellerDashboard(auth, page, size, status)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

// GET /dashboard/admin?page=1&size=20&status=paid
func (h *DashboardHandler) GetAdminDashboard(w http.ResponseWriter, r *http.Request) error {
	page, size := parsePagination(r)
	status := r.URL.Query().Get("status")

	res, err := h.dashboardService.GetAdminDashboard(page, size, status)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func parsePagination(r *http.Request) (page, size int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	size, _ = strconv.Atoi(r.URL.Query().Get("size"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return
}
