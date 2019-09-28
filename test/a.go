package a

import (
	"fmt"

	"github.com/nakario/findqueries/test/b"
)

func callHoge() {
	f := hoge()
	fmt.Println(f.GetPiyo())
}

func hoge() *b.Fuga {
	return b.NewFuga()
}
