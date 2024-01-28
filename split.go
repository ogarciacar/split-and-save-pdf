package split

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/split-and-save-pdf", new(PDF))
}

type PDF struct{}

func (*PDF) SplitAndSave(pdfData []byte, filename string) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	log.Printf("Splitting %s\n", filename)

	err := os.WriteFile(filename, pdfData, 0644)
	if err != nil {
		log.Printf("error: %s\n", err)
		os.Exit(1)
	}

	outputFilePrefix := fmt.Sprintf("%s_", strings.TrimSuffix(filename, ".pdf"))
	cmd := exec.Command("pdftk", fmt.Sprintf("./%s", filename), "burst", "output", outputFilePrefix)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error: %s\n", err)
	}

	log.Printf("Split completed %s %s\n", filename, out)
}

// example
// func main() {

// 	filename := "JRV_2046.pdf"

// 	pdfFile, err := os.Open(filename)
// 	if err != nil {
// 		os.Exit(1)
// 	}
// 	defer pdfFile.Close()

// 	fileData, err := io.ReadAll(pdfFile)
// 	if err != nil {
// 		os.Exit(1)
// 	}

// 	pdf := PDF{}
// 	pdf.SplitAndSave(fileData, "filename.pdf")
// }
