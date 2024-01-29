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

func (*PDF) Save(pdfData []byte, filename string) error {

	removeTimestampFromLogs()

	log.Printf("Save PDF %s\n", filename)

	err := os.WriteFile(filename, pdfData, 0644)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	}

	return nil
}

func (*PDF) Split(filename string) error {

	removeTimestampFromLogs()

	log.Printf("Split PDF %s\n", filename)

	outputFilePrefix := fmt.Sprintf("%s_", strings.TrimSuffix(filename, ".pdf"))
	cmd := exec.Command("pdftk", fmt.Sprintf("./%s", filename), "burst", "output", outputFilePrefix)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	}

	log.Printf("PDF split done %s %s\n", filename, out)
	return nil
}

func (*PDF) Read(filename string) []byte {

	removeTimestampFromLogs()

	log.Printf("Read PDF %s\n", filename)
	pdfFile, err := os.Open(filename)
	if err != nil {
		log.Printf("error: %s\n", err)
		return make([]byte, 0)
	}
	defer pdfFile.Close()

	fileData, err := io.ReadAll(pdfFile)

	if err != nil {
		log.Printf("error: %s\n", err)
		return make([]byte, 0)
	}

	return fileData
}
