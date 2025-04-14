package http

import pkgErrors "reup/pkg/errors"

var (
	errMissBody             = pkgErrors.NewHTTPError(3000, "Invalid request payload")
	errWrongPaginationQuery = pkgErrors.NewHTTPError(3100, "Wrong pagination query")
)
