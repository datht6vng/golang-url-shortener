package util

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/belinskiydm/bsconv"
	"github.com/jxskiss/base62"
)

func GenerateShortLink(ID string) string {
	newShortUrl, err := bsconv.ConvertFromDec(ID, 62)
	if err != nil {
		panic(err)
	}
	return newShortUrl + SignUrl(newShortUrl)
}
func SignUrl(input string) string {
	secret := "Secret key: !@#$%^&*()432144adsfdsafhk12312532dfhsajkfghjkg732478321er1234hkjdhf78234h"
	hashSecret := hmac.New(sha256.New, []byte(secret))
	hashSecret.Write([]byte(input))
	return base62.EncodeToString(hashSecret.Sum(nil))[:5]
}

// var (
// 	// CharacterSet consists of 62 characters [0-9][A-Z][a-z].
// 	Base         = int64(62)
// 	CharacterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
// )

// // Encode returns a base62 representation as
// // string of the given integer number.
// func Encode(num int64) string {
// 	b := make([]byte, 0)
// 	// loop as long the num is bigger than zero
// 	for num > 0 {
// 		// receive the rest
// 		r := math.Mod(float64(num), float64(Base))

// 		// devide by Base
// 		num /= Base

// 		// append chars
// 		b = append([]byte{CharacterSet[int(r)]}, b...)
// 	}
// 	return string(b)
// }

// // Decode returns a integer number of a base62 encoded string.
// func Decode(s string) (int64, error) {
// 	var r, pow int64

// 	// loop through the input
// 	for i, v := range s {
// 		// convert position to power
// 		pow = int64(len(s) - (i + 1))

// 		// IndexRune returns -1 if v is not part of CharacterSet.
// 		pos := int64(strings.IndexRune(CharacterSet, v))

// 		if pos == -1 {
// 			return pos, errors.New("invalid character: " + string(v))
// 		}

// 		// calculate
// 		r += pos * int64(math.Pow(float64(Base), float64(pow)))
// 	}
// 	return int64(r), nil
// }
// func TrimTimeStamp(s int64, base float64) int64 {
// 	return s % int64(math.Pow(10, base))
// }
