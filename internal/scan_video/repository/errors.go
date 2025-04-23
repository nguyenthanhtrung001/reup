package repository

import (
	pkgErrors "github.com/nguyenthanhtrung001/reup/pkg/errors"
)

var (
	ErrDocumentNotFound = pkgErrors.NewHTTPError(6105, "not found document")
)
