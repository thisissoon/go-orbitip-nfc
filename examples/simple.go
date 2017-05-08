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
		orbitip.DefaultRoot,
		orbitip.DefaultExt,
		orbitip.Handlers{
			orbitip.CardReadCmd: func(rv orbitip.ResponseValues, p orbitip.Params) error {
				fmt.Println(fmt.Sprintf("NFC read from %s", p.UID))
				rv.Beep(orbitip.LongBeep)
				rv.Ext(orbitip.HTML)
				rv.Deny()
				return nil
			},
		})
	go srv.ListenAndServe()
	defer srv.Shutdown(context.Background())
	C := make(chan os.Signal, 1)
	signal.Notify(C, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-C
}
