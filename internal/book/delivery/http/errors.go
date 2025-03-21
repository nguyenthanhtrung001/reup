package http

import pkgErrors "book-store/pkg/errors"

var (
	errMissBody             = pkgErrors.NewHTTPError(3000, "Invalid request payload")
	errWrongPaginationQuery = pkgErrors.NewHTTPError(3100, "Wrong pagination query")
)
