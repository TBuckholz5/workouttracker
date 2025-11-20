package hash

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hash string, password string) error
}

type BcryptHash struct{}

func NewBcryptHasher() *BcryptHash {
	return &BcryptHash{}
}

func (b *BcryptHash) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (b *BcryptHash) VerifyPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
