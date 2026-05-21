package order

import (
	"github.com/devlucas-java/klyp-shop/pkg/pagination"
)

type OrdersPageResponse struct {
	Pagination pagination.OutPutPagination `json:"pagination"`
	Items      []*OrderResponse            `json:"items"`
}
