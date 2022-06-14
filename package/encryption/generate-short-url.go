package encryption

import (
	"fmt"
	"github.com/belinskiydm/bsconv"
)

func GenerateShortLink(ID int64) string {
	newShortURL, err := bsconv.ConvertFromDec(fmt.Sprintf("%v", ID), 62)
	if err != nil {
		panic(err)
	}
	return newShortURL + Signature(newShortURL)[:4]
}
