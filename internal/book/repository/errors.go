package repository

import (
	pkgErrors "reup/pkg/errors"
)

var (
	ErrDocumentNotFound = pkgErrors.NewHTTPError(6105, "not found document")
)
