package repository

import (
	pkgErrors "book-store/pkg/errors"
)

var (
	ErrDocumentNotFound = pkgErrors.NewHTTPError(6105, "not found document")
)
