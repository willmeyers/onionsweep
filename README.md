<p align="center">
  <img src="./onionsweep.png" alt="onionsweep" width="200"/>
</p>

# Onionsweep

A fast way to check for live .onion URLs.

## Installation

After cloning this repo, run

```
make
```

`onionsweep` will be built in the current directory.

## Usage (CLI)

Pass a newline separated list of .onion URLs into `onionsweep`. For example,

```bash
cat onions.txt | ./onionsweep -workers 3 -torlistenaddr 127.0.0.1:9050 -timeout 10 > results.txt
```

Results tab separated list of URLs, status codes, and response times, document title, and lengths

```
http://l7....onion  403 live
http://vx....onion  403 live
http://me....onion  0   dead   Get "http://me....onion": socks connect tcp 127.0.0.1:9050->me....onion:80: unknown error host unreachable
http://wc....onion  403 live
http://ne....onion  200 live
http://to....onion  200 live
http://to....onion  200 live
```

### Useful CLI Flags

```
Usage of onionsweep:
  -timeout int
        Timeout in seconds (default 10)
  -torListenAddr string
        Tor listen address (default "127.0.0.1:9050")
  -workers int
        Number of workers to use (default 10)
```
