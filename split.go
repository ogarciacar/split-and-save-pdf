package split

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
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

func (*PDF) Split(filename string) int {

	removeTimestampFromLogs()

	log.Printf("Split PDF %s\n", filename)

	outputFilePrefix := fmt.Sprintf("%s_", strings.TrimSuffix(filename, ".pdf"))

	cmd := exec.Command("pdftk", fmt.Sprintf("./%s", filename), "burst", "output", outputFilePrefix)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error: %s\n", err)
		return 0
	}

	log.Printf("PDF split done %s %s\n", filename, out)

	numberOfPagesCmdStr := fmt.Sprintf("pdftk %s dump_data | grep NumberOfPages | awk -F ': ' '{print $2}' | tr -d '\n'", filename)

	out, err = exec.Command("bash", "-c", numberOfPagesCmdStr).CombinedOutput()
	if err != nil {
		log.Printf("error: %s\n", err)
		return 0
	}

	byteToInt, err := strconv.Atoi(string(out))
	if err != nil {
		log.Printf("error: %s\n", err)
		return 0
	}

	log.Printf("%s has %d pages\n", filename, byteToInt)

	return byteToInt
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
