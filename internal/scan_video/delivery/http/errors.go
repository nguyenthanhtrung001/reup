package http

import pkgErrors "github.com/nguyenthanhtrung001/reup/pkg/errors"

var (
	errMissBody             = pkgErrors.NewHTTPError(3000, "Invalid request payload")
	errWrongPaginationQuery = pkgErrors.NewHTTPError(3100, "Wrong pagination query")
	errDoneAllProxyScan     = pkgErrors.NewHTTPError(3200, "Error done all proxy scan")
	errDoneProxyScan        = pkgErrors.NewHTTPError(3201, "Error done proxy scan")
	errGetAllProxyScan      = pkgErrors.NewHTTPError(3202, "Error get all proxy scan")
	errGetProxyScanRandom   = pkgErrors.NewHTTPError(3203, "Error get proxy scan random")
	errRequiredQueryParam   = pkgErrors.NewHTTPError(3204, "Required query param: proxy_ip")
)
