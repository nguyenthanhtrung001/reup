package http

import (
	"book-store/internal/models"
	"book-store/pkg/jwt"

	pkgErrors "book-store/pkg/errors"

	"github.com/gin-gonic/gin"
)

func (h handler) authenticateUser(c *gin.Context) (models.Scope, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "control.http.authenticateUser.GetPayloadFromContext: unauthorized")
		return models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	sc := jwt.NewScope(payload)
	return sc, nil
}

func (h handler) processBookReq(c *gin.Context) (bookRequest, error) {
	ctx := c.Request.Context()
	var rqBody bookRequest
	if err := c.ShouldBindJSON(&rqBody); err != nil {
		h.l.Errorf(ctx, "book.http.create.shouldBindJson")
		return bookRequest{}, err
	}
	return rqBody, nil
}

func (h handler) processListBookReq(c *gin.Context) (listBookReq, error) {
	ctx := c.Request.Context()
	var rqBody listBookReq
	if err := c.ShouldBindQuery(&rqBody); err != nil {
		h.l.Errorf(ctx, "book.http.List.ShouldBindQuery: %v", err)
		return listBookReq{}, err
	}
	return rqBody, nil
}
