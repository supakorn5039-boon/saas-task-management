package security

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// bcryptCost defaults to 14 (slow on purpose — that's the whole point of
// bcrypt). Tests in CI override this to bcrypt.MinCost (4) via the
// BCRYPT_COST env var, which makes each Register/Login take ~1ms instead of
// ~1s. Production must never set BCRYPT_COST below ~10.
var bcryptCost = func() int {
	if v, ok := os.LookupEnv("BCRYPT_COST"); ok {
		if n, err := strconv.Atoi(v); err == nil && n >= bcrypt.MinCost && n <= bcrypt.MaxCost {
			return n
		}
	}
	return 14
}()

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
