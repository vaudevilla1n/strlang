package quote

import (
	"fmt"
	"rsc.io/quote"
)

func Quote() {
	fmt.Println(quote.Go());
	fmt.Println(quote.Glass());
	fmt.Println(quote.Opt());
	fmt.Println(quote.Hello());
}
