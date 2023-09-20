package main

import (
	"flag"

	onionSweep "github.com/willmeyers/onionsweep"
)

func main() {
	workers := flag.Int("workers", 10, "Number of workers to use")
	timeout := flag.Int("timeout", 10, "Timeout in seconds")
	torListenAddr := flag.String("torListenAddr", "127.0.0.1:9050", "Tor listen address")

	flag.Parse()

	config := &onionSweep.OnionSweepConfig{
		Workers:          *workers,
		TimeOutInSeconds: *timeout,
		TorListenAddr:    *torListenAddr,
	}

	onionSweep.Run(config)
}
