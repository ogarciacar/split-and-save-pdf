package split

import (
	"fmt"
	"io"
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

func removeTimestampFromLogs() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

func (*PDF) Save(pdfData []byte, filename string) {

	removeTimestampFromLogs()

	log.Printf("Saving PDF as %s\n", filename)

	err := os.WriteFile(filename, pdfData, 0644)
	if err != nil {
		log.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func (*PDF) Split(filename string) {

	removeTimestampFromLogs()

	log.Printf("Splitting %s into pages\n", filename)

	outputFilePrefix := fmt.Sprintf("%s_", strings.TrimSuffix(filename, ".pdf"))
	cmd := exec.Command("pdftk", fmt.Sprintf("./%s", filename), "burst", "output", outputFilePrefix)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error: %s\n", err)
	}

	log.Printf("PDF split completed %s %s\n", filename, out)
}

func (*PDF) Read(filename string) []byte {

	removeTimestampFromLogs()

	log.Printf("Reading PDF bytes from %s\n", filename)
	pdfFile, err := os.Open(filename)
	if err != nil {
		log.Printf("error: %s\n", err)
		os.Exit(1)
	}
	defer pdfFile.Close()

	fileData, err := io.ReadAll(pdfFile)

	if err != nil {
		log.Printf("error: %s\n", err)
		os.Exit(1)
	}

	return fileData
}
