# ethspam

`ethspam` generates an infinite stream of realistic read-only Ethereum JSONRPC queries,
anchored around the latest block with some amount of random jitter. The latest state is updated every 15 seconds, so it can run continuously without becoming stale.

Per second, ethspam generates around 500,000 lines, or 120 megabytes of data, on a modern portable laptop.

Also makes for an okay superniche screensaver.


## Getting Started

A few options:

- [Grab a binary release](https://github.com/shazow/ethspam/releases), or
- Build from source: `go get github.com/shazow/ethspam`, or
- Use the Docker Hub image: [`docker pull shazow/ethspam`](https://hub.docker.com/r/shazow/ethspam)

Then run it and enjoy. Ethspam will emit output but it can be throttled by backflow pressure from the consumer. If the process you're piping to isn't consuming the output fast enough, ethspam will slow down.

```
$ ethspam | head
...
```


## License

MIT
