package onionsweep

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type Worker struct {
	ID         int
	OnionSweep *OnionSweep
}

// Cache of already checked hosts
var cache = make(map[string]bool)
var cacheLock = &sync.Mutex{}

func NewWorker(id int, onionSweep *OnionSweep) *Worker {
	return &Worker{
		ID:         id,
		OnionSweep: onionSweep,
	}
}

func (w *Worker) Run() {
	defer w.OnionSweep.Wg.Done()

	for job := range w.OnionSweep.Jobs {
		parsedUrl, err := url.Parse(job)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing URL: %s\n", err.Error())
			continue
		}

		cacheLock.Lock()
		resolvable, ok := cache[parsedUrl.Host]
		cacheLock.Unlock()
		if ok && !resolvable {
			w.writeResult(job, "cached", 0, "dead")
			continue
		}

		req, err := http.NewRequest("GET", job, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing URL: %s\n", err.Error())
			continue
		}

		resp, err := w.OnionSweep.Client.Do(req)
		if err != nil {
			cacheLock.Lock()
			cache[parsedUrl.Host] = false
			cacheLock.Unlock()
			w.writeResult(job, err.Error(), 0, "dead")
			continue
		}

		cacheLock.Lock()
		cache[parsedUrl.Host] = true
		cacheLock.Unlock()

		w.writeResult(job, "", resp.StatusCode, "live")
	}
}

func (w *Worker) writeResult(url, errorReason string, statusCode int, liveOrDead string) {
	w.OnionSweep.Results <- fmt.Sprintf("%s\t%d\t%s\t%s", url, statusCode, liveOrDead, errorReason)
}
