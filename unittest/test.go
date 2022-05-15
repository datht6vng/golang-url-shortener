package unittest

import (
	"fmt"
	"server_go/util"
)

func TestShortener() {
	links := []string{
		"http://localhost",
		"https://www.google.com",
		"https://chat.zalo.me",
		"https://www.facebook.com",
		"https://shopee.vn/?gclid=Cj0KCQjwyYKUBhDJARIsAMj9lkFkSqvjjEhRWG5JJMszI5Z30EbDCERc1oZ91tA1RKAh7U5CkzPd5k8aAk6-EALw_wcB"}
	for _, link := range links {
		fmt.Println(util.GenerateShortLink(link))
	}
}
