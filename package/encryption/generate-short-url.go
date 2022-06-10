package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"trueid-shorten-link/config"

	"github.com/belinskiydm/bsconv"
	"github.com/jxskiss/base62"
)

func GenerateShortLink(ID string) string {
	newShortUrl, err := bsconv.ConvertFromDec(ID, 62)
	if err != nil {
		panic(err)
	}
	return newShortUrl + Sign(newShortUrl)
}
func Sign(input string) string {
	secret := config.Config.Key.Secret
	hashSecret := hmac.New(sha256.New, []byte(secret))
	hashSecret.Write([]byte(input))
	return base62.EncodeToString(hashSecret.Sum(nil))[:4]
}
