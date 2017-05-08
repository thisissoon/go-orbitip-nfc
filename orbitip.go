package orbitip

import (
	"fmt"
	"net/http"
	"strings"
)

// Default path and file ext the reader sends requests too: /orbit.php
var (
	DefaultRoot = "/orbit"
	DefaultExt  = PHP
)

// TODO: Document
type BeepDuration int

// TODO: Document
var (
	ShortBeep = BeepDuration(1)
	LongBeep  = BeepDuration(0)
)

// The ResponseValues type is a map of NFC response commands and their values
type ResponseValues map[string]string

// HeartBeatInterval will set the HB response command setting the readers
// heart beat interval in seconds, value must be between 1 and 9999
func (rv ResponseValues) HeartBeatInterval(seconds int) error {
	if seconds < 1 || seconds > 9999 {
		return fmt.Errorf("HB can only be between 1 and 9999: %d", seconds)
	}
	rv["HB"] = fmt.Sprintf("%04d", seconds)
	return nil
}

// Beep returns the BEEP response command triggering the reader to beep for a
// long or short duration
func (rv ResponseValues) Beep(bd BeepDuration) {
	rv["BEEP"] = fmt.Sprintf("%d", bd)
}

// TODO: Implementation
// Clock returns the CK response command setting the clock on the reader.
// The year must be between 2000 and 2099.
func (rv ResponseValues) Clock(year, month, day, hour, min, sec int) error {
	return nil
}

// TODO: Implementation
// ClockCalibration returns the CCAL response command.
func (rv ResponseValues) ClockCalibration(opperand string, value string) error {
	return nil
}

// Grant returns the GRNT response command, setting Set reader to a “grant access” state.
// The orange LED will be set to ON for seconds provided. If the relay is set to “active” then the
// coil will be powered to engage the relay from NO state to NC state for the seconds provided
// and return to NO state.
func (rv ResponseValues) Grant(seconds int) error {
	if seconds < 1 || seconds > 99 {
		return fmt.Errorf("GRNT can only be between 1 and 99: %d", seconds)
	}
	rv["GRNT"] = fmt.Sprintf("%02d", seconds)
	return nil
}

// Deny returns the DENY response command setting the reader into a deny state
func (rv ResponseValues) Deny() {
	rv["DENY"] = ""
}

// Root returns the ROOT response command, allowing customisation of the web server
// root path the reader calls, this can be up to 8 characters. Use 000000000 to reset
// to default configuration.
func (rv ResponseValues) Root(value string) error {
	if value != "000000000" || len(value) > 8 {
		return fmt.Errorf("ROOT can be no more than 8 characters: %s is %d", value, len(value))
	}
	rv["ROOT"] = value
	return nil
}

// Ext returns the EXT response command, configuring the web server path extenstion used
// by the reader for HTTP calls
func (rv ResponseValues) Ext(ext Ext) error {
	rv["EXT"] = fmt.Sprintf("%d", ext.ID)
	return nil
}

// Command defines NFC reader command types
type Command string

// String returns the command string
func (c Command) String() string {
	return string(c)
}

// NFC Commands
var (
	PowerUpCmd     = Command("PU") // Power Up
	HeartBeatCmd   = Command("HB") // Heartbeat
	CardReadCmd    = Command("CO") // Card Opperation
	LevelChangeCmd = Command("SW") // Level Change
	PingCmd        = Command("PG") // Ping
)

// Ext defines NFC path extenstion (.php, .asp etc)
type Ext struct {
	ID   uint8
	Name string
}

// Return the string of the extenstion
func (e Ext) String() string {
	return e.Name
}

// Supported extenstions by the reader
var (
	PHP  = Ext{0, ".php"} // Default
	ASP  = Ext{1, ".asp"}
	CFM  = Ext{2, ".cfm"}
	PL   = Ext{3, ".pl"}
	HTM  = Ext{4, ".htm"}
	HTML = Ext{5, ".html"}
	ASPX = Ext{6, ".aspx"}
	JSP  = Ext{7, ".jsp"}
)

// Params sent by NFC reader
type Params struct {
	Date     string
	Time     string
	ID       string
	ULen     string
	UID      string
	Command  string
	Version  string
	Contact1 string
	Contact2 string
	SID      string
	Data     string
	PSRC     string
	MD5      string
	MAC      string
	Relay    string
	SD       string
}

// The HandlerFunc type allows the use of ordinary functions as NFC command handlers.
type HandlerFunc func(response ResponseValues, parms Params) error

// Handlers maps a NFC command to a specific HandlerFunc.
type Handlers map[Command]HandlerFunc

// Set sets a HandlerFunc for the given NFC command
func (h Handlers) Set(cmd Command, fn HandlerFunc) {
	h[cmd] = fn
}

// Del removes a HandleFunc for the given NFC command
func (h Handlers) Del(cmd Command) {
	delete(h, cmd)
}

// The ServeMux type implements the http.Handler interface for serving
// NFC requests. This can be passed directly into a http.Server instance on the
// Handler property.
type ServeMux struct {
	handlers Handlers
}

// Handlers returns the NFC Command to HandlerFunc map
func (m *ServeMux) Handlers() Handlers {
	return m.handlers
}

// ServeHTTP handles an NFC HTTP request, first checking if a HandleFunc
// exists for the NFC command and if so calling that function. Any data returned
// by the HandleFunc is written to the http response.
func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fn, ok := m.handlers[Command(values.Get("cmd"))]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	rv := make(ResponseValues)
	err := fn(rv, Params{
		Date:     values.Get("date"),
		Time:     values.Get("time"),
		ID:       values.Get("id"),
		ULen:     values.Get("ulen"),
		UID:      values.Get("uid"),
		Command:  values.Get("cmd"),
		Version:  values.Get("ver"),
		Contact1: values.Get("contact1"),
		Contact2: values.Get("contact2"),
		SID:      values.Get("sid"),
		Data:     values.Get("data"),
		PSRC:     values.Get("psrc"),
		MD5:      values.Get("md5"),
		MAC:      values.Get("mac"),
		Relay:    values.Get("relay"),
		SD:       values.Get("sd"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(rv) > 0 {
		lines := []string{"<ORBIT>"}
		for k, v := range rv {
			line := fmt.Sprintf("%s=%s", k, v)
			if v == "" {
				line = k
			}
			lines = append(lines, line)
		}
		lines = append(lines, "</ORBIT>")
		w.Write([]byte(strings.Join(lines, "\n")))
	}
}

// NewServeMux constructs a new ServeMux
func NewServeMux(handlers Handlers) *ServeMux {
	return &ServeMux{
		handlers: handlers,
	}
}

// New constructs a new HTTP server for running the NFC HTTP server
func New(addr string, root string, ext Ext, handlers Handlers) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("%s%s", root, ext), NewServeMux(handlers))
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
