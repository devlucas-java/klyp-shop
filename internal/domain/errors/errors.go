package errors

var (
	ErrUserNotFound = &AppError{
		Code:    "USER_NOT_FOUND",
		Message: "user not found",
		Status:  404,
	}
	ErrUserAlreadyExists = &AppError{
		Code:    "USER_ALREADY_EXISTS",
		Message: "user already exists",
		Status:  409,
	}
	ErrInvalidCredentials = &AppError{
		Code:    "INVALID_CREDENTIALS",
		Message: "invalid credentials",
		Status:  401,
	}
	ErrUserIsAlreadySeller = &AppError{
		Code:    "USER_ALREADY_SELLER",
		Message: "user is already a seller",
		Status:  409,
	}
	ErrSellerCannotBuy = &AppError{
		Code:    "SELLER_CANNOT_BUY",
		Message: "sellers cannot place orders",
		Status:  403,
	}
	ErrInvalidRole = &AppError{
		Code:    "INVALID_ROLE",
		Message: "role transition is not allowed",
		Status:  400,
	}

	ErrSellerNotFound = &AppError{
		Code:    "SELLER_NOT_FOUND",
		Message: "seller not found",
		Status:  404,
	}
	ErrSellerNotOwner = &AppError{
		Code:    "SELLER_NOT_OWNER",
		Message: "seller is not the owner of this product",
		Status:  403,
	}

	ErrProductNotFound = &AppError{
		Code:    "PRODUCT_NOT_FOUND",
		Message: "product not found",
		Status:  404,
	}
	ErrInsufficientStock = &AppError{
		Code:    "INSUFFICIENT_STOCK",
		Message: "insufficient stock",
		Status:  422,
	}
	ErrInvalidProductPrice = &AppError{
		Code:    "INVALID_PRODUCT_PRICE",
		Message: "product price must be greater than zero",
		Status:  400,
	}

	ErrOrderNotFound = &AppError{
		Code:    "ORDER_NOT_FOUND",
		Message: "order not found",
		Status:  404,
	}
	ErrOrderCannotBeCancelled = &AppError{
		Code:    "ORDER_CANNOT_BE_CANCELLED",
		Message: "order cannot be cancelled in current status",
		Status:  422,
	}
	ErrOrderNotPending = &AppError{
		Code:    "ORDER_NOT_PENDING",
		Message: "order is not pending",
		Status:  422,
	}

	ErrPaymentNotFound = &AppError{
		Code:    "PAYMENT_NOT_FOUND",
		Message: "payment not found",
		Status:  404,
	}
	ErrPaymentAlreadyDone = &AppError{
		Code:    "PAYMENT_ALREADY_CONFIRMED",
		Message: "payment already confirmed",
		Status:  409,
	}

	ErrReviewNotFound = &AppError{
		Code:    "REVIEW_NOT_FOUND",
		Message: "review not found",
		Status:  404,
	}
	ErrReviewAlreadyExists = &AppError{
		Code:    "REVIEW_ALREADY_EXISTS",
		Message: "user already reviewed this product",
		Status:  409,
	}
	ErrInvalidRating = &AppError{
		Code:    "INVALID_RATING",
		Message: "rating must be between 1 and 5",
		Status:  400,
	}

	ErrInvalidCategoryName = &AppError{
		Code:    "INVALID_CATEGORY_NAME",
		Message: "category name must not contain special characters",
		Status:  400,
	}
	ErrCategoryNotFound = &AppError{
		Code:    "CATEGORY_NOT_FOUND",
		Message: "category not found",
		Status:  404,
	}
)
