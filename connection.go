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

// Is there a better way to do this?
// I manually calculated the size in C
//
const cupsDestTSize = 32
const cupsOptionTSize = 16

var enumDestsChan chan *Dest

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

		options := make(map[string]string)
		for j := 0; j < int(dest.num_options)-1; j++ {
			var opt *C.cups_option_t
			opt = dest.options
			options[C.GoString(opt.name)] = C.GoString(opt.value)
		}
		d.options = options

		destsArr = append(destsArr, d)
	}

	// free the pointers
	C.cupsFreeDests(destCount, dests)

	conn.NumDests = goDestCount
	conn.Dests = destsArr
}

// TODO: implement cupsEnumDests()
