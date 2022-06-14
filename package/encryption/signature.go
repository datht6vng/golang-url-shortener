package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"trueid-shorten-link/package/config"

	"github.com/jxskiss/base62"
)

func Signature(input string) string {
	secret := config.Config.Key.Secret
	hashSecret := hmac.New(sha256.New, []byte(secret))
	hashSecret.Write([]byte(input))
	return base62.EncodeToString(hashSecret.Sum(nil))
}
