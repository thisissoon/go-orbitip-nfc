/*
This package offers is a simple HTTP server for Gemini Orbit IP NFC readers.

Overview

The Orbit IP attempts to HTTP requests to a webserver, by default
runs at 192.168.7.191 on port 80, serving path on /orbit.php. The reader sends
several URL query parameters which contain meta data about the request, for
example the command it is performing, data, date time and so on.

This package allows you to construct and run a simple HTTP server for serving
these requests and returning responses to the reader.

A simple implementation looks like this:

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
			":80",
			orbitip.DefaultRoot,
			orbitip.DefaultExt,
			orbitip.Handlers{
				orbitip.PowerUpCmd: func(rv orbitip.ResponseValues, p orbitip.Params) error {
					fmt.Println("Power Up")
					rv.UI(orbitip.UI{RedFlash: true, BuzzerIntermittent: true}, 5, 50)
					return nil
				},
				orbitip.CardReadCmd: func(rv orbitip.ResponseValues, p orbitip.Params) error {
					fmt.Println("Card Read")
					rv.UI(orbitip.UI{GreenFlash: true, BuzzerIntermittent: true}, 3, 50)
					return nil
				},
				orbitip.PingCmd: func(rv orbitip.ResponseValues, p orbitip.Params) error {
					fmt.Println("Ping")
					return nil
				},
				orbitip.LevelChangeCmd: func(rv orbitip.ResponseValues, p orbitip.Params) error {
					fmt.Println("Level Change", p.Contact1, p.Contact2)
					return nil
				},
				orbitip.LevelChangeCmd: func(rv orbitip.ResponseValues, p orbitip.Params) error {
					fmt.Println("Heart Beat")
					return nil
				},
			})
		go srv.ListenAndServe()
		defer srv.Shutdown(context.Background())
		C := make(chan os.Signal, 1)
		signal.Notify(C, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-C
	}
*/
package orbitip
