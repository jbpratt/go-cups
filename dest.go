package cups

/*
#cgo LDFLAGS: -lcups
#include "cups/cups.h"
*/
import "C"
import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const printerStateIdle = "3"
const printerStatePrinting = "4"
const printerStateStopped = "5"

// Dest is a struct for each printer
type Dest struct {
	Name      string
	IsDefault bool
	options   map[string]string
	// internal keys for status
}

// PrintFile prints a file
// TODO: Complete implementation
func (d *Dest) PrintFile(filename, mimeType string) (int, error) {
	// check mime type
	mimeTypes, err := GetMimeTypes("/usr/share/cups/mime/mime.types")
	if err != nil {
		log.Fatal(err)
	}

	i := sort.SearchStrings(mimeTypes, mimeType)
	if i < len(mimeTypes) && mimeTypes[i] == mimeType {
		fmt.Println("Valid")
	}
	// check file
	id := C.cupsPrintFile(C.CString(d.Name), C.CString(filename),
		C.CString("Test Print"), C.int(len(d.options)), nil)

	if id == 0 {
		return -1, errors.New(fmt.Sprintf("Failed to print with: %d error code", C.cupsLastError()))
	}
	return int(id), nil
}

// Status returns the status of the dest
func (d *Dest) Status() string {
	var returnMessage string

	// Return status of dest
	value, ok := d.options["printer-state"]

	if ok != true {
		returnMessage = "printer state key does not exist"
	}

	switch value {
	case printerStateIdle:
		returnMessage = "idle"
		break
	case printerStatePrinting:
		returnMessage = "printing"
		break
	case printerStateStopped:
		returnMessage = "stopped"
		break
	default:
		returnMessage = "error"
		break
	}

	return returnMessage

}

// GetOption returns the options
// TODO: Complete implementation
// func (d *Dest) GetOption(keyName string) string {
// 	// Return option
// 	return ""
// }

// GetOptions returns a map of the dest options
func (d *Dest) GetOptions() map[string]string {
	// Return option
	return d.options
}

// TestPrint prints a test page
func (d *Dest) TestPrint() int {
	var numOptions C.int
	var options *C.cups_option_t
	var jobID C.int

	// resolve path/to/test/file
	// same across linux/osx; TODO: bsd
	test := "/usr/share/cups/data/testfile"
	// Print a single file
	jobID = C.cupsPrintFile(C.CString(d.Name), C.CString(test),
		C.CString("Test Print"), numOptions, options)

	return int(jobID)
}

// TODO: Implement CancelJob
// cupsCancelDestJob

// GetMimeTypes returns a slice of strings
func GetMimeTypes(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var mimeTypes []string

	for _, l := range lines {
		mt := strings.Fields(l)
		if len(mt) > 0 {
			if !strings.Contains(mt[0], "#") && !strings.Contains(mt[0], "(") {
				mimeTypes = append(mimeTypes, mt[0])
			}
		}
	}

	return mimeTypes, nil
}
