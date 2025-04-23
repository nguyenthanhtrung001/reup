package http

import (
	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/internal/user/usecase"
	pkgErrors "github.com/nguyenthanhtrung001/reup/pkg/errors"
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/paginator"
	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary List
// @Schemes
// @Description List
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Produce json
// @Tags User - Admin
// @Accept json
// @Produce json
// @Success 200 {object} listResponse
// @Failure 400 {object} response.Resp "Bad Request,Error..."
// @Router /api/v1/admin/users/list [GET]
func (h handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processListRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.http.List.processListRequest: %v", err)
		response.Error(c, err)
		return
	}

	var pagQuery paginator.PaginatorQuery
	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		h.l.Errorf(ctx, "user.http.Get.gin.ShouldBindQuery: %v", err)
		response.Error(c, errWrongPaginationQuery)
		return
	}

	pagQuery.Adjust()

	o, err := h.userUC.List(ctx, sc, usecase.ListInput{
		PaginatorQuery: pagQuery,
		Filter: usecase.ListInputFilter{
			Name: req.Name,
		},
	})
	if err != nil {
		h.l.Errorf(ctx, "user.http.List.userUC.List: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, newListResponse(o))
}

// @Summary Detail Admin
// @Schemes
// @Description Detail Admin
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Param id path string true "id"
// @Produce json
// @Tags User - Admin
// @Accept json
// @Produce json
// @Success 200 {object} detailAdminResponse
// @Failure 400 {object} response.Resp "Bad Request,Error..."
// @Router /api/v1/admin/users/{id} [GET]
func (h handler) DetailAdmin(c *gin.Context) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "user.http.List.GetPayloadFromContext: unauthorized")
		response.Error(c, pkgErrors.NewUnauthorizedHTTPError())
	}

	sc := jwt.NewScope(payload)

	id := c.Param("id")

	o, err := h.userUC.GetUserById(ctx, sc, id)
	if err != nil {
		h.l.Errorf(ctx, "user.http.DetailAdmin.userUC.DetailAdmin: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, newDetailAdminResponse(o))
}

func (h handler) processUpdateAdminRequest(c *gin.Context) (updateAdminRequest, models.Scope, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "user.http.processUpdateAdminRequest.GetPayloadFromContext: unauthorized")
		return updateAdminRequest{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	sc := jwt.NewScope(payload)

	var rqBody updateAdminRequest
	if err := c.ShouldBindJSON(&rqBody); err != nil {
		h.l.Errorf(ctx, "user.http.processUpdateAdminRequest.ShouldBindJSON: %v", errInvalidBody)
		return updateAdminRequest{}, models.Scope{}, err
	}

	return rqBody, sc, nil
}

// @Summary Update By Admin
// @Schemes
// @Description Update By Admin
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Param id path string true "id"
// @Param body body updateAdminRequest true "Update By Admin "
// @Produce json
// @Tags User - Admin
// @Accept json
// @Produce json
// @Success 200 {object} detailAdminResponse
// @Failure 400 {object} response.Resp "Bad Request,Error..."
// @Router /api/v1/admin/users/{id} [PUT]
func (h handler) UpdateByAdmin(c *gin.Context) {
	ctx := c.Request.Context()

	rqBody, sc, err := h.processUpdateAdminRequest(c)
	if err != nil {
		response.AdminError(c, err)
		return
	}

	id := c.Param("id")

	input := usecase.UpdateByAdminInput{
		Email:    rqBody.Email,
		FullName: rqBody.FullName,
		Phone:    rqBody.Phone,
		GroupID:  rqBody.GroupID,
		Discount: rqBody.Discount,
		Priority: rqBody.Priority,
	}

	o, err := h.userUC.UpdateByAdmin(ctx, sc, input, id)
	if err != nil {
		h.l.Errorf(ctx, "user.http.UpdateAdmin.userUC.UpdateByAdmin: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, newDetailAdminResponse(o))
}

// @Summary Create Api Key
// @Schemes
// @Description Create Api Key
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Param id query string true "id"
// @Produce json
// @Tags User - Admin
// @Accept json
// @Produce json
// @Success 200 {object} detailAdminResponse
// @Failure 400 {object} response.Resp "Bad Request,Error..."
// @Router /api/v1/admin/users/{id}/api-key [POST]
func (h handler) CreateApiKey(c *gin.Context) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "user.http.processCreateApiKeyRequest.GetPayloadFromContext: unauthorized")
		response.Error(c, pkgErrors.NewUnauthorizedHTTPError())
		return
	}

	sc := jwt.NewScope(payload)

	id := c.Param("id")

	o, err := h.userUC.CreateApiKeyUser(ctx, sc, id)
	if err != nil {
		h.l.Errorf(ctx, "user.http.CreateApiKey.userUC.CreateApiKey: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, newDetailAdminResponse(o))
}
