package onionsweep

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

type OnionSweepConfig struct {
	Workers          int
	TimeOutInSeconds int
	TorListenAddr    string
}

type OnionSweep struct {
	Config  *OnionSweepConfig
	Jobs    chan string
	Results chan string
	Client  *http.Client
	Wg      sync.WaitGroup
}

func (onionSweep *OnionSweep) newHttpClient() (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", onionSweep.Config.TorListenAddr, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URL: %s\n", err.Error())
		return nil, err
	}

	httpTransport := &http.Transport{}
	httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	httpClient := &http.Client{Transport: httpTransport, Timeout: time.Duration(onionSweep.Config.TimeOutInSeconds) * time.Second}

	return httpClient, nil
}

func Run(onionSweepConfig *OnionSweepConfig) {
	onionSweep := &OnionSweep{
		Config:  onionSweepConfig,
		Jobs:    make(chan string),
		Results: make(chan string),
	}

	httpClient, err := onionSweep.newHttpClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	onionSweep.Client = httpClient

	onionSweep.Wg.Add(onionSweepConfig.Workers)
	for i := 0; i < onionSweepConfig.Workers; i++ {
		worker := NewWorker(i, onionSweep)
		go worker.Run()
	}

	go func() {
		onionSweep.Wg.Wait()
		close(onionSweep.Results)
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			job := scanner.Text()
			onionSweep.Jobs <- job
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading standard input: %s", err)
		}

		close(onionSweep.Jobs)
	}()

	for result := range onionSweep.Results {
		fmt.Println(result)
	}
}
