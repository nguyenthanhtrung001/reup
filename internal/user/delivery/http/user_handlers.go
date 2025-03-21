package http

import (
	"book-store/internal/models"
	"book-store/internal/user/usecase"
	pkgErrors "book-store/pkg/errors"
	"book-store/pkg/jwt"
	"book-store/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h handler) processRegisterRequest(c *gin.Context) (registerRequest, error) {
	ctx := c.Request.Context()

	var rqBody registerRequest

	if err := c.ShouldBindJSON(&rqBody); err != nil {
		h.l.Errorf(ctx, "user.http.Register.ShouldBindJSON: %v", errInvalidBody)
		return registerRequest{}, errMissingParams
	}
	if !rqBody.validateEmail() {
		h.l.Errorf(ctx, "user.http.Register.validateEmail: %v", errInvalidEmail)
		return registerRequest{}, errInvalidEmail
	}

	if !rqBody.validatePassword() {
		h.l.Errorf(ctx, "user.http.Register.validatePassword: %v", errInvalidPassword)
		return registerRequest{}, errInvalidPassword
	}

	if !rqBody.validatePhone() {
		h.l.Errorf(ctx, "user.http.Register.validatePhone: %v", errInvalidPhone)
		return registerRequest{}, errInvalidPhone
	}
	return rqBody, nil
}

func (h handler) processLoginRequest(c *gin.Context) (loginRequest, error) {
	ctx := c.Request.Context()

	var rqBody loginRequest
	if err := c.ShouldBindJSON(&rqBody); err != nil {
		h.l.Errorf(ctx, "user.http.Login.ShouldBindJSON: %v", errInvalidBody)
		return loginRequest{}, errInvalidBody
	}

	return rqBody, nil
}

func (h handler) processProfileRequest(c *gin.Context) (profileRequest, models.Scope, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "user.http.Profile.GetPayloadFromContext: unauthorized")
		return profileRequest{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	sc := jwt.NewScope(payload)

	return profileRequest{}, sc, nil
}

func (h handler) processUpdateProfileRequest(c *gin.Context) (updateProfileRequest, models.Scope, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "user.http.UpdateProfile.GetPayloadFromContext: unauthorized")
		return updateProfileRequest{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	sc := jwt.NewScope(payload)

	var rqBody updateProfileRequest
	if err := c.ShouldBindJSON(&rqBody); err != nil {
		h.l.Errorf(ctx, "user.http.UpdateProfile.ShouldBindJSON: %v", errInvalidBody)
		return updateProfileRequest{}, models.Scope{}, errInvalidBody
	}

	return rqBody, sc, nil
}

func (h handler) processUpdatePasswordRequest(c *gin.Context) (updatePasswordRequest, models.Scope, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "user.http.UpdatePassword.GetPayloadFromContext: unauthorized")
		return updatePasswordRequest{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	sc := jwt.NewScope(payload)

	var rqBody updatePasswordRequest
	if err := c.ShouldBindJSON(&rqBody); err != nil {
		h.l.Errorf(ctx, "user.http.UpdatePassword.ShouldBindJSON: %v", errInvalidBody)
		return updatePasswordRequest{}, models.Scope{}, errInvalidBody
	}

	return rqBody, sc, nil
}

func (h handler) processListRequest(c *gin.Context) (listRequest, models.Scope, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "user.http.List.GetPayloadFromContext: unauthorized")
		return listRequest{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	sc := jwt.NewScope(payload)

	var req listRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "user.http.List.ShouldBindQuery: %v", err)
		return listRequest{}, sc, err
	}

	return req, sc, nil
}

// @Summary Register
// @Schemes
// @Description Register
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Param body body registerRequest true "Register "
// @Produce json
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} adminItemResponse
// @Failure 400 {object} response.Resp "Bad Request,Error User already registered(4003)"
// @Router /api/v1/users/register [POST]
func (h handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	rqBody, err := h.processRegisterRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.http.Register.processRegisterRequest: %v", err)
		response.Error(c, err)
		return
	}
	input := usecase.RegisterInput{
		FullName: rqBody.FullName,
		Email:    rqBody.Email,
		Phone:    rqBody.Phone,
		Password: rqBody.Password,
	}

	// Add referer_id
	if rqBody.RefererID != "" {
		input.RefererID = rqBody.RefererID
	}

	if err := h.userUC.Register(ctx, input); err != nil {
		if err == usecase.ErrUserAlreadyExist {
			h.l.Errorf(ctx, "user.http.Register.userUC.Register: %v", errUserFound)
			response.Error(c, errUserFound)
			return
		}
		h.l.Errorf(ctx, "user.http.Register.userUC.Register: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}

// @Summary Login
// @Schemes
// @Description Login
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Param body body loginRequest true "Login "
// @Produce json
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} loginResponse
// @Failure 400 {object} response.Resp "Bad Request,User not found(4005)"
// @Router /api/v1/users/login [POST]
func (h handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	rqBody, err := h.processLoginRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.http.Login.processLoginRequest: %v", err)
		response.Error(c, err)
		return
	}

	input := usecase.LoginInput{
		Email:    rqBody.Email,
		Password: rqBody.Password,
	}
	o, err := h.userUC.Login(ctx, input)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			h.l.Errorf(ctx, "user.http.Login.userUC.Login: %v", errUserNotFound)
			response.Error(c, errUserNotFound)
			return
		}
		if err == usecase.ErrInvalidValidation {
			h.l.Errorf(ctx, "user.http.Login.userUC.Login: %v", errInvalidValidation)
			response.Error(c, errInvalidValidation)
			return
		}
		h.l.Errorf(ctx, "user.http.Login.userUC.Login: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, newLoginResponse(o))
}

// @Summary Detail Profile
// @Schemes
// @Description Detail Profile
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Produce json
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} profileResponse
// @Failure 400 {object} response.Resp "Bad Request,Error..."
// @Router /api/v1/users/profile [GET]
func (h handler) Profile(c *gin.Context) {
	ctx := c.Request.Context()

	_, sc, err := h.processProfileRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.http.Profile.processProfileRequest: %v", err)
		response.Error(c, err)
		return
	}

	o, err := h.userUC.Profile(ctx, sc)
	if err != nil {
		h.l.Errorf(ctx, "user.http.Profile.userUC.Profile: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, newProfileResponse(o))
}

// @Summary Update Profile
// @Schemes
// @Description Update Profile
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Param body body updateProfileRequest true "Update Profile"
// @Produce json
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.Resp
// @Failure 400 {object} response.Resp "Bad Request,Error..."
// @Router /api/v1/users/profile [PUT]
func (h handler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	rqBody, sc, err := h.processUpdateProfileRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.http.UpdateProfile.processUpdateProfileRequest: %v", err)
		response.Error(c, err)
		return
	}

	input := usecase.UpdateProfileInput{
		Name: rqBody.FullName,
	}
	err = h.userUC.UpdateProfile(ctx, sc, input)
	if err != nil {
		h.l.Errorf(ctx, "user.http.UpdateProfile.userUC.UpdateProfile: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}

// @Summary Update Password
// @Schemes
// @Description Update Password
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default(*)
// @Param User-Agent header string false "User-Agent" default(Swagger-Codegen/1.0.0/go)
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjYyYzc3MzU1MmRmMzZjMGJkMjUxNDdkIiwiZ3JvdXBfaWQiOiI2NGY0NDUyNTlkODNkM2JkZDg0ZGZjOWEiLCJncm91cF9yb2xlIjoiYWRtaW4iLCJleHBpcmVkX2F0IjoiMjAyNC0wNS0yN1QxMDo1Nzo0Ny40Mjg1NTgrMDc6MDAiLCJleHAiOjE3MTY3ODIyNjd9.pb6sLIq4F2FDXE2ASWacYRzI5qs1ae48_DeQ2b3jJLU)"
// @Param Language header string false "Language" default(en)
// @Param body body updatePasswordRequest true "Update Password"
// @Produce json
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.Resp
// @Failure 400 {object} response.Resp "Bad Request,Invalid validation(6002)."
// @Router /api/v1/users/password [PUT]
func (h handler) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()

	rqBody, sc, err := h.processUpdatePasswordRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.http.UpdatePassword.processUpdatePasswordRequest: %v", err)
		response.Error(c, err)
		return
	}

	input := usecase.UpdatePasswordInput{
		OldPassword: rqBody.OldPassword,
		NewPassword: rqBody.NewPassword,
	}

	if err := h.userUC.UpdatePassword(ctx, sc, input); err != nil {
		if err == usecase.ErrInvalidValidation {
			h.l.Errorf(ctx, "user.http.UpdatePassword.userUC.UpdatePassword: %v", errInvalidValidation)
			response.Error(c, errInvalidValidation)
			return
		}
		h.l.Errorf(ctx, "user.http.UpdatePassword.userUC.UpdatePassword: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}
