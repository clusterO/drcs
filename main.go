package main

import (
	"DCRS/dcrs"
	network "DCRS/net"
	"time"
)

func main() {
	go func() {
		network.Listen()
	}()

	time.Sleep(1 * time.Second)
	dcrs.Cli()
}
