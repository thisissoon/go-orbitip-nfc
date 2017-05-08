# Go Orbit IP NFC

A Go library for running an NFC reader server for Gemini 2000 Orbit IP devices (http://www.gemini2k.com/orbit-ip-poe-nfc-smart-card-reader/).

## Example

``` go
package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	orbitip "github.com/thisissoon/go-orbitip-nfc"
)

func main() {
	handlers := make(orbitip.Handlers)
	handlers.Set(orbitip.CO, func(v url.Values) ([]byte, error) {
		fmt.Println("NFC Read")
		return nil, nil
	})
	srv := orbitip.New(
		":8000",
		orbitip.DEFAULT_ROOT,
		orbitip.DEFAULT_EXT,
		handlers)
	go srv.ListenAndServe()
	defer srv.Shutdown(context.Background())
	C := make(chan os.Signal, 1)
	signal.Notify(C, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-C
}
```
