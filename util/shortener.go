package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/itchyny/base58-go"
)

// func sha256Of(input string) []byte {
// 	algorithm := sha256.New()
// 	algorithm.Write([]byte(input))
// 	return algorithm.Sum(nil)
// }

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

// func GenerateShortLink(initialLink string) string {
// 	u, err := gouuid.NewV4()
// 	if err != nil {
// 		panic(err)
// 	}
// 	urlHashBytes := sha256Of(initialLink + u.String())
// 	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
// 	shortUrl := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))[:8]
// 	return shortUrl + SignUrl(shortUrl)
// }

func GenerateShortLink(ID int64) string {
	newShortUrl := Encode(ID)
	return newShortUrl + SignUrl(newShortUrl)
}
func SignUrl(input string) string {
	secret := "Secret key: !@#$%^&*()432144adsfdsafhk12312532dfhsajkfghjkg732478321er1234hkjdhf78234h"
	hashSecret := hmac.New(sha256.New, []byte(secret))
	hashSecret.Write([]byte(input))
	return base58Encoded(hashSecret.Sum(nil))[:4]
}

var (
	// CharacterSet consists of 62 characters [0-9][A-Z][a-z].
	Base         = int64(62)
	CharacterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Encode returns a base62 representation as
// string of the given integer number.
func Encode(num int64) string {
	b := make([]byte, 0)

	// loop as long the num is bigger than zero
	for num > 0 {
		// receive the rest
		r := math.Mod(float64(num), float64(Base))

		// devide by Base
		num /= Base

		// append chars
		b = append([]byte{CharacterSet[int(r)]}, b...)
	}

	return string(b)
}

// Decode returns a integer number of a base62 encoded string.
func Decode(s string) (int64, error) {
	var r, pow int64

	// loop through the input
	for i, v := range s {
		// convert position to power
		pow = int64(len(s) - (i + 1))

		// IndexRune returns -1 if v is not part of CharacterSet.
		pos := int64(strings.IndexRune(CharacterSet, v))

		if pos == -1 {
			return pos, errors.New("invalid character: " + string(v))
		}

		// calculate
		r += pos * int64(math.Pow(float64(Base), float64(pow)))
	}
	return int64(r), nil
}
