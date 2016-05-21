package main

import (
	//"fmt"
	"github.com/m4tta/goddis/goddis"
)

func main() {
	goddis := goddis.NewGoddis()
	goddis.Listen(6379)
}
