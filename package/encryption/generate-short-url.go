package encryption

import (
	"github.com/belinskiydm/bsconv"
)

func GenerateShortLink(ID string) string {
	newShortURL, err := bsconv.ConvertFromDec(ID, 62)
	if err != nil {
		panic(err)
	}
	return newShortURL + Signature(newShortURL)[:4]
}
