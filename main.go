package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/INFURA/go-ethlibs/node"
)

// Top query ratios

/*
      3 "eth_accounts"
      4 "eth_getStorageAt"
      4 "eth_syncing"
      7 "net_peerCount"
     12 "net_listening"
     14 "eth_gasPrice"
     16 "eth_sendRawTransaction"
     25 "net_version"
     30 "eth_getTransactionByBlockNumberAndIndex"
     38 "eth_getBlockByHash"
     45 "eth_estimateGas"
     88 "eth_getCode"
    252 "eth_getLogs"
    255 "eth_getTransactionByHash"
    333 "eth_blockNumber"
    390 "eth_getTransactionCount"
    399 "eth_getBlockByNumber"
    545 "eth_getBalance"
    607 "eth_getTransactionReceipt"
   1928 "eth_call"
*/

var defaultWeb3Endpoint = "https://mainnet.infura.io/v3/af500e495f2d4e7cbcae36d0bfa66bcb" // Versus API key on Infura

func exit(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(code)
}

func main() {
	gen := generator{}
	installDefaults(&gen)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := node.NewClient(ctx, defaultWeb3Endpoint)
	if err != nil {
		exit(1, "failed to make a new client: %s", err)
	}
	mkState := stateProducer{
		client: client,
	}

	// stateChannel 😂
	stateChannel := make(chan State, 1)

	// We don't need a high quality randomness source, just for benchmark shuffling
	randSrc := rand.NewSource(time.Now().UnixNano())
	go func() {
		state := liveState{
			idGen:   &idGenerator{},
			randSrc: randSrc,
		}
		for {
			newState, err := mkState.Refresh(&state)
			if err != nil {
				exit(2, "failed to refresh state")
			}
			select {
			case stateChannel <- newState:
			case <-ctx.Done():
				return
			}

			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
			}
		}
	}()

	state := <-stateChannel
	for {
		// Update state when a new one is emitted
		select {
		case state = <-stateChannel:
		case <-ctx.Done():
			return
		default:
		}
		if err := gen.Query(os.Stdout, state); err != nil {
			exit(2, "failed to generate query: %s", err)
		}
	}
}

type Generator func(io.Writer, State) error

type RandomQuery struct {
	Method   string
	Weight   int64
	Generate Generator
}

type generator struct {
	queries     []RandomQuery // sorted by weight asc
	totalWeight int64
}

// Add inserts a random query generator with a weighted probability. Not
// goroutine-safe, should be run once during initialization.
func (g *generator) Add(query RandomQuery) {
	if g.queries == nil {
		g.queries = make([]RandomQuery, 1)
	} else {
		g.queries = append(g.queries, RandomQuery{})
	}
	// Maintain weight sort
	idx := sort.Search(len(g.queries), func(i int) bool { return g.queries[i].Weight < query.Weight })
	copy(g.queries[idx+1:], g.queries[idx:])
	g.queries[idx] = query
	g.totalWeight += query.Weight
}

// Query selects a generator based on proportonal weighted probability and
// writes the query from the generator.
func (g *generator) Query(w io.Writer, s State) error {
	if len(g.queries) == 0 {
		return errors.New("no query generators available")
	}

	weight := s.RandInt64() % g.totalWeight

	var current int64
	for _, q := range g.queries {
		// TODO: Test for off-by-one
		current += q.Weight
		if current >= weight {
			return q.Generate(w, s)
		}
	}

	panic("off by one bug in weighted query selection")
}
