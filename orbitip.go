package orbitip

import (
	"fmt"
	"net/http"
	"net/url"
)

// Default path and file ext the reader sends requests too: /orbit.php
var (
	DEFAULT_ROOT = "/orbit"
	DEFAULT_EXT  = PHP
)

// Type for definging NFC reader commands
type Command string

// Support stringer interface
func (c Command) String() string {
	return string(c)
}

// NFC Commands
var (
	PU = Command("PU") // Power Up
	HB = Command("HB") // Heartbeat
	CO = Command("CO") // Card Opperation
	SW = Command("SW") // Level Change
	PG = Command("PG") // Ping
)

// NFC Path Extenstion
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

// NFC command handler function
type HandlerFunc func(parms url.Values) ([]byte, error)

// Mapping of NFC commands to handler functions
type Handlers map[Command]HandlerFunc

// Sets the handler function for a given NFC command
func (h Handlers) Set(cmd Command, fn HandlerFunc) {
	h[cmd] = fn
}

// Removes a handler function for the given NFC command
func (h Handlers) Del(cmd Command) {
	delete(h, cmd)
}

// NFC Server Mux, implementing the http.Handler interface
// This can be passed directly into a http.Server instance on the Handler property
type ServeMux struct {
	handlers Handlers
}

// Implements the http.Handler interface
func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fn, ok := m.handlers[Command(values.Get("cmd"))]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := fn(values)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

// Constructs a new NFC Server Mux
func NewServeMux(handlers Handlers) *ServeMux {
	return &ServeMux{
		handlers: handlers,
	}
}

// Constructs a new ready to use HTTP server for running the NFC HTTP server
func New(addr string, root string, ext Ext, handlers Handlers) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("%s%s", root, ext), NewServeMux(handlers))
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
