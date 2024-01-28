package split

import (
	"log"
	"os"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/split-and-save-pdf", new(PDF))
}

type PDF struct{}

func (*PDF) SplitAndSave(pdfData []byte, filename string) {
	log.Println("Split and save")

	savePDFToFile(pdfData, filename)
}

func savePDFToFile(pdfData []byte, filename string) error {
	err := os.WriteFile(filename, pdfData, 0644)
	if err != nil {
		return err
	}
	return nil
}
