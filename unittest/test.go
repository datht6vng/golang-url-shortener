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
		"https://shopee.vn/Bao-Tay-Ch%C6%A1i-Game-ff-Pubg-Li%C3%AAn-Qu%C3%A2n...-G%C4%83ng-tay-ch%C6%A1i-game-Ch%E1%BB%91ng-M%E1%BB%93-H%C3%B4i-Si%C3%AAu-Nh%E1%BA%A1y-Co-Gi%C3%A3n-C%E1%BB%B1c-T%E1%BB%91t-B%E1%BA%A3o-H%C3%A0nh-12-Th%C3%A1ng-i.390877573.10033483585?sp_atk=bd8d9e3f-4459-4800-81ce-894c15a44921&xptdk=bd8d9e3f-4459-4800-81ce-894c15a44921",
		"https://shopee.vn/Thi%E1%BA%BFt-B%E1%BB%8B-%C4%90i%E1%BB%87n-T%E1%BB%AD-cat.11036132",
		"https://shopee.vn/",
		"http://localhost",
		"https://www.google.com",
		"https://www.google.com",
		"https://chat.zalo.me",
		"https://www.facebook.com",
		"https://shopee.vn/?gclid=Cj0KCQjwyYKUBhDJARIsAMj9lkFkSqvjjEhRWG5JJMszI5Z30EbDCERc1oZ91tA1RKAh7U5CkzPd5k8aAk6-EALw_wcB",
		"https://www.google.com/search?q=url-shortener&ei=kxiBYp69BIKzmAW5n5iwBg&ved=0ahUKEwje0ujJ5eH3AhWCGaYKHbkPBmYQ4dUDCA4&uact=5&oq=url-shortener&gs_lcp=Cgdnd3Mtd2l6EAMyCAgAEAcQHhATMgQIABATMgQIABATMgYIABAeEBMyBggAEB4QEzIGCAAQHhATMgYIABAeEBMyBggAEB4QEzIGCAAQHhATMgYIABAeEBM6BwgAEEcQsAM6BwgAELADEEM6EgguEMcBENEDEMgDELADEEMYAToGCAAQBxAeOggIABANEB4QE0oECEEYAEoECEYYAFCuF1iQG2DkH2gCcAF4AIABWogB0wKSAQE0mAEAoAEByAELwAEB2gEECAEYCA&sclient=gws-wiz",
		"https://www.google.com/search?q=news&ei=OKiBYrX5EqiMseMPq8Cb2A0&ved=0ahUKEwi19dTI7uL3AhUoRmwGHSvgBtsQ4dUDCA4&uact=5&oq=news&gs_lcp=Cgdnd3Mtd2l6EAMyCggAELEDEIMBEEMyBQgAEIAEMgUIABCABDIFCAAQgAQyCAgAEIAEELEDMgUIABCABDIICAAQgAQQsQMyBAgAEEMyBQgAEIAEMgsIABCABBCxAxCDAToHCAAQRxCwAzoHCAAQsAMQQzoSCC4QxwEQ0QMQyAMQsAMQQxgBOgsILhCABBCxAxCDAToOCC4QgAQQsQMQxwEQ0QM6CwguEIAEEMcBEKMCOgUILhCABDoNCC4QsQMQxwEQ0QMQQzoECC4QQzoRCC4QgAQQsQMQgwEQxwEQowI6CAguELEDEIMBOgoILhCxAxDUAhBDOgsILhCABBCxAxDUAjoKCAAQsQMQgwEQCkoECEEYAEoECEYYAFC1CFj5CmCVDGgBcAF4AIABe4gBjAOSAQMzLjGYAQCgAQHIAQvAAQHaAQQIARgI&sclient=gws-wiz"}
	for _, link := range links {
		url := util.GenerateShortLink(link)
		signature := util.SignUrl(url)
		fmt.Println("----------------------------------------------------")
		fmt.Println(url)
		fmt.Println(signature)
	}
}
