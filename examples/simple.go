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
