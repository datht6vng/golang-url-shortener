package unittest

import (
	"fmt"
	"net/url"
	"server_go/util"
)

func TestValidUrl() {
	links := []string{
		"www.gooog.com",
		"http://localhost",
		"https://www.google.com",
		"https://chat.zalo.me",
		"https://www.facebook.com",
		"https://shopee.vn/?gclid=Cj0KCQjwyYKUBhDJARIsAMj9lkFkSqvjjEhRWG5JJMszI5Z30EbDCERc1oZ91tA1RKAh7U5CkzPd5k8aAk6-EALw_wcB",
		"https://www.google.com/search?q=url-shortener&ei=kxiBYp69BIKzmAW5n5iwBg&ved=0ahUKEwje0ujJ5eH3AhWCGaYKHbkPBmYQ4dUDCA4&uact=5&oq=url-shortener&gs_lcp=Cgdnd3Mtd2l6EAMyCAgAEAcQHhATMgQIABATMgQIABATMgYIABAeEBMyBggAEB4QEzIGCAAQHhATMgYIABAeEBMyBggAEB4QEzIGCAAQHhATMgYIABAeEBM6BwgAEEcQsAM6BwgAELADEEM6EgguEMcBENEDEMgDELADEEMYAToGCAAQBxAeOggIABANEB4QE0oECEEYAEoECEYYAFCuF1iQG2DkH2gCcAF4AIABWogB0wKSAQE0mAEAoAEByAELwAEB2gEECAEYCA&sclient=gws-wiz",
		"https://www.google.com/search?q=news&ei=OKiBYrX5EqiMseMPq8Cb2A0&ved=0ahUKEwi19dTI7uL3AhUoRmwGHSvgBtsQ4dUDCA4&uact=5&oq=news&gs_lcp=Cgdnd3Mtd2l6EAMyCggAELEDEIMBEEMyBQgAEIAEMgUIABCABDIFCAAQgAQyCAgAEIAEELEDMgUIABCABDIICAAQgAQQsQMyBAgAEEMyBQgAEIAEMgsIABCABBCxAxCDAToHCAAQRxCwAzoHCAAQsAMQQzoSCC4QxwEQ0QMQyAMQsAMQQxgBOgsILhCABBCxAxCDAToOCC4QgAQQsQMQxwEQ0QM6CwguEIAEEMcBEKMCOgUILhCABDoNCC4QsQMQxwEQ0QMQQzoECC4QQzoRCC4QgAQQsQMQgwEQxwEQowI6CAguELEDEIMBOgoILhCxAxDUAhBDOgsILhCABBCxAxDUAjoKCAAQsQMQgwEQCkoECEEYAEoECEYYAFC1CFj5CmCVDGgBcAF4AIABe4gBjAOSAQMzLjGYAQCgAQHIAQvAAQHaAQQIARgI&sclient=gws-wiz"}
	for _, link := range links {
		_, err := url.ParseRequestURI(link)
		if err != nil {
			fmt.Println("Not valid")
		} else {
			fmt.Println("Valid")
		}
	}
}
func TestShortener() {
	links := []string{
		"http://localhost",
		"https://www.google.com",
		"https://chat.zalo.me",
		"https://www.facebook.com",
		"https://shopee.vn/?gclid=Cj0KCQjwyYKUBhDJARIsAMj9lkFkSqvjjEhRWG5JJMszI5Z30EbDCERc1oZ91tA1RKAh7U5CkzPd5k8aAk6-EALw_wcB",
		"https://www.google.com/search?q=url-shortener&ei=kxiBYp69BIKzmAW5n5iwBg&ved=0ahUKEwje0ujJ5eH3AhWCGaYKHbkPBmYQ4dUDCA4&uact=5&oq=url-shortener&gs_lcp=Cgdnd3Mtd2l6EAMyCAgAEAcQHhATMgQIABATMgQIABATMgYIABAeEBMyBggAEB4QEzIGCAAQHhATMgYIABAeEBMyBggAEB4QEzIGCAAQHhATMgYIABAeEBM6BwgAEEcQsAM6BwgAELADEEM6EgguEMcBENEDEMgDELADEEMYAToGCAAQBxAeOggIABANEB4QE0oECEEYAEoECEYYAFCuF1iQG2DkH2gCcAF4AIABWogB0wKSAQE0mAEAoAEByAELwAEB2gEECAEYCA&sclient=gws-wiz",
		"https://www.google.com/search?q=news&ei=OKiBYrX5EqiMseMPq8Cb2A0&ved=0ahUKEwi19dTI7uL3AhUoRmwGHSvgBtsQ4dUDCA4&uact=5&oq=news&gs_lcp=Cgdnd3Mtd2l6EAMyCggAELEDEIMBEEMyBQgAEIAEMgUIABCABDIFCAAQgAQyCAgAEIAEELEDMgUIABCABDIICAAQgAQQsQMyBAgAEEMyBQgAEIAEMgsIABCABBCxAxCDAToHCAAQRxCwAzoHCAAQsAMQQzoSCC4QxwEQ0QMQyAMQsAMQQxgBOgsILhCABBCxAxCDAToOCC4QgAQQsQMQxwEQ0QM6CwguEIAEEMcBEKMCOgUILhCABDoNCC4QsQMQxwEQ0QMQQzoECC4QQzoRCC4QgAQQsQMQgwEQxwEQowI6CAguELEDEIMBOgoILhCxAxDUAhBDOgsILhCABBCxAxDUAjoKCAAQsQMQgwEQCkoECEEYAEoECEYYAFC1CFj5CmCVDGgBcAF4AIABe4gBjAOSAQMzLjGYAQCgAQHIAQvAAQHaAQQIARgI&sclient=gws-wiz"}
	for _, link := range links {
		fmt.Println(util.GenerateShortLink(link))
	}
}
