package basic

import (
	"fmt"
	"rsc.io/quote"
)

func QuotePrint() {
	fmt.Printf("\nQuote ::: %s", quote.Go())
}
