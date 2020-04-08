# ethspam

`ethspam` generates an infinite stream of realistic read-only Ethereum JSONRPC queries,
anchored around the latest block with some amount of random jitter. The latest state is updated every 15 seconds, so it can run continuously without becoming stale.

Per second, ethspam generates around 500,000 lines, or 120 megabytes of data, on a modern portable laptop.

Also makes for an okay superniche screensaver.


## Setup

[Grab a binary release](https://github.com/shazow/ethspam/releases) or build from source.


## Build docker image and push to docker hub

```
docker build -t shazow/ethspam .
docker push shazow/ethspam
```
## License

MIT
