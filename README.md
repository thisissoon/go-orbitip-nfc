# Go Orbit IP NFC

A Go library for running an NFC reader server for Gemini 2000 Orbit IP devices (http://www.gemini2k.com/orbit-ip-poe-nfc-smart-card-reader/).

## Example

``` go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	orbitip "github.com/thisissoon/go-orbitip-nfc"
)

func main() {
	srv := orbitip.New(
		":8000",
		orbitip.DEFAULT_ROOT,
		orbitip.DEFAULT_EXT,
		orbitip.Handlers{
			orbitip.CO: func(p orbitip.Params) ([]byte, error) {
				fmt.Println(fmt.Sprintf("NFC read from %s", p.UID))
				return nil, nil
			},
		})
	go srv.ListenAndServe()
	defer srv.Shutdown(context.Background())
	C := make(chan os.Signal, 1)
	signal.Notify(C, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-C
}
```
