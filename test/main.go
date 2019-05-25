package main

/*
#cgo LDFLAGS: -lcups
#include "cups/cups.h"
*/
import "C"
import "fmt"

func main() {

	id := C.cupsPrintFile(C.CString("ZJ-58"), C.CString("../test.txt"), C.CString("test print"), 0, nil)

	if id == 0 {
		fmt.Printf("%+v\n", C.cupsLastError())
	}
}
