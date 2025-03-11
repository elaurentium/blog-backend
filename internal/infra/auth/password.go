package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

type PasswordService struct {
	config PasswordConfig
}

func NewPasswordService() *PasswordService {
	return &PasswordService{
		config: PasswordConfig{
			Time:    3,
			Memory:  64 * 1024,
			Threads: 4,
			KeyLen:  32,
			SaltLen: 16,
		},
	}
}

func (s *PasswordService) HashPassword(password string) (string, string, error) {
	salt := make([]byte, s.config.SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", "", err
	}

	hash := argon2.IDKey([]byte(password), salt, s.config.Time, s.config.Memory, s.config.Threads, s.config.KeyLen)

	saltStr := base64.RawStdEncoding.EncodeToString(salt)
	hashStr := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, s.config.Memory, s.config.Time, s.config.Threads, saltStr, hashStr)

	return encodedHash, saltStr, nil
}

func (s *PasswordService) VerifyPassword(password, encodedHash, saltStr string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	hashToCompare := argon2.IDKey([]byte(password), salt, s.config.Time, s.config.Memory, s.config.Threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1
}