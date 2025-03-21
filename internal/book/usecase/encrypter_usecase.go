package usecase

import "context"

func (uc implUseCase) Create(ctx context.Context, text string) (string, error) {
	return uc.encrypter.Encrypt(text)
}
