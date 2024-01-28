package split

import (
	"log"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/split-and-save-pdf", new(PDF))
}

type PDF struct{}

func (*PDF) SplitAndSave() {
	log.Println("Split and save")
}
