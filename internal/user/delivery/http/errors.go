package http

import (
	pkgErrors "github.com/nguyenthanhtrung001/reup/pkg/errors"
)

var (
	errInvalidPassword = pkgErrors.NewHTTPError(1400, "Invalid password")
	errInvalidEmail    = pkgErrors.NewHTTPError(4002, "Invalid email")
	errUserFound       = pkgErrors.NewHTTPError(4003, "User already registered")
	errUserNotFound    = pkgErrors.NewHTTPError(4004, "User not found")
	errMissingParams   = pkgErrors.NewHTTPError(4005, "Missing params, please check email, phone, fullname, password")
	errInvalidPhone    = pkgErrors.NewHTTPError(4006, "Invalid phone")

	errWrongPaginationQuery = pkgErrors.NewHTTPError(6000, "Wrong pagination query")
	errInvalidBody          = pkgErrors.NewHTTPError(6001, "Invalid body")
	errInvalidValidation    = pkgErrors.NewHTTPError(6002, "Invalid validation")
)
