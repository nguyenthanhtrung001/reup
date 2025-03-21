package encrypter

type Encrypter interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type implEncrypter struct {
	key string
}

func NewEncrypter(key string) Encrypter {
	return &implEncrypter{
		key: key,
	}
}
