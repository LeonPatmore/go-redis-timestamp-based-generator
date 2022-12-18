package example

import (
	"fmt"

	"github.com/leonpatmore/gotemplate/internal/simple"
)

func Handle() {
	fmt.Print(simple.Add(1, 3))
}
