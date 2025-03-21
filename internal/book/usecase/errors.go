package usecase

import pkgErrors "reup/pkg/errors"

var (
	ErrPackageConfirmed = pkgErrors.NewHTTPError(1000, "package already confirmed")
	ErrInvalidID        = pkgErrors.NewHTTPError(2000, "ID is invalid")
	ErrNotFound         = pkgErrors.NewHTTPError(2001, "resource not found")
)
