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
// TODO: should we swap from taking in mimetype and instead
//	detecting content type through http.DetectContentType
// leak
func (d *Dest) PrintFile(filename, mimeType string) (int, error) {
	// check mime type
	mimeTypes, err := GetMimeTypes("/usr/share/cups/mime/mime.types")
	if err != nil {
		log.Fatal(err)
	}

	sort.Strings(mimeTypes)
	i := sort.SearchStrings(mimeTypes, mimeType)
	if i < len(mimeTypes) && mimeTypes[i] == mimeType {

		// TODO: check if file exists

		status := d.Status()
		if status != "idle" {
			log.Fatal("printer is not idle")
		}

		// print file
		id := C.cupsPrintFile(C.CString(d.Name), C.CString(filename),
			C.CString("Test Print"), C.int(len(d.options)), nil)

		if id == 0 {
			return 1, errors.New(
				fmt.Sprintf("failed to print with error code: %d %s",
					C.cupsLastError(), C.GoString(C.cupsLastErrorString())))
		}
		return int(id), nil
	}
	return 1, errors.New("invalid mime type")
}

// TODO: PrintFiles
// cupsPrintFiles

// Status returns the status of the dest
// TODO: return string, error
func (d *Dest) Status() string {
	var status string

	// Return status of dest
	value, ok := d.options["printer-state"]

	if ok != true {
		status = "printer state key does not exist"
	}

	switch value {
	case printerStateIdle:
		status = "idle"
		break
	case printerStatePrinting:
		status = "printing"
		break
	case printerStateStopped:
		status = "stopped"
		break
	default:
		status = "error"
		break
	}

	return status
}

// GetOption returns the options
func (d *Dest) GetOption(keyName string) (string, error) {
	options := d.GetOptions()
	if v, found := options[keyName]; found {
		return v, nil
	}
	return "", errors.New("option not found")
}

// GetOptions returns a map of the dest options
func (d *Dest) GetOptions() map[string]string {
	// Return option
	return d.options
}

// TODO: CreateJob
// cupsCreateJob
// cupsStartDocument
// cupsWriteRequestData
// cupsFinishDocument

// TODO: GetJobs
// cupsGetJobs
func (d *Dest) GetJobs() {

	// check dest status
	// get jobs
}

// CancelJob takes in a jobID and cancels the job with
// that ID. TODO: implement CUPS_JOBID_CURRENT & CUPS_JOBID_ALL
func (d *Dest) CancelJob(jobID int) error {
	if jobID == 0 {
		return errors.New("supply job ID TODO: implement CUPS_JOBID_CURRENT")
	}

	id := C.cupsCancelJob(C.CString(d.Name), C.int(jobID))

	if id == 0 {
		return errors.New(
			fmt.Sprintf("failed to cancel job with error code: %d %v",
				C.cupsLastError(), C.GoString(C.cupsLastErrorString())))
	}
	return nil
}

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
