package usecase

import "errors"

var (
	ErrPackageConfirmed  = errors.New("package already confirmed")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidValidation = errors.New("invalid validation")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrInvalidAuthType   = errors.New("invalid auth type")
	ErrInvalidAuth       = errors.New("error max update auth")
	ErrInvalidRotate     = errors.New("error max update rotate")
	ErrNumberRotate      = errors.New("error number rotate")
	ErrMaxAuthIP         = errors.New("error max auth ip")
	ErrInvalidIP         = errors.New("error invalid ip")

	ErrTrialPackageCreated = errors.New("trial package already created")
	ErrInvalidPackageName  = errors.New("invalid package name")
	ErrThisFunctionStopped = errors.New("can not use this function on this version")
	ErrOrderRefuned        = errors.New("order already refunded")
	// Topup
	ErrTopupConfirmed = errors.New("topup already confirmed")
)
