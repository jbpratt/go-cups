package cups

/*
#cgo LDFLAGS: -lcups
#include "cups/cups.h"
*/
import "C"

// Connection is a struct containing information about a CUPS connection
type Connection struct {
	isDefault bool
	address   string
	NumDests  int
	Dests     []*Dest
}

// Refresh updates the list of destinations and their state
func (c *Connection) Refresh() {
	updateDefaultConnection(c)
}

// NewDefaultConnection creates a new default CUPS connection
func NewDefaultConnection() *Connection {

	connection := &Connection{
		isDefault: true,
	}
	updateDefaultConnection(connection)

	return connection
}

func updateDefaultConnection(conn *Connection) {
	var dests *C.cups_dest_t
	destCount := C.cupsGetDests(&dests)
	goDestCount := int(destCount)

	var destsArr []*Dest

	for i := 0; i < goDestCount; i++ {

		dest := dests
		d := &Dest{
			Name: C.GoString(dest.name),
		}

		options := []string{
			"auth-info-required",
			"printer-info",
			"printer-is-accepting-jobs",
			"printer-is-shared",
			"printer-location",
			"printer-make-and-model",
			"printer-state",
			"printer-state-change-time",
			"printer-state-reasons",
			"printer-type",
		}
		res := make(map[string]string, len(options))

		for _, option := range options {
			o := C.cupsGetOption(C.CString(option), dest.num_options, dest.options)
			res[option] = C.GoString(o)
		}

		d.options = res

		destsArr = append(destsArr, d)
	}

	// free the pointers
	C.cupsFreeDests(destCount, dests)

	conn.NumDests = goDestCount
	conn.Dests = destsArr
}

// TODO: implement cupsEnumDests()
