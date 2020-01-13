package main

import (
	"errors"
	"io"
	"sort"
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

func main() {

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

// Add inserts a random query generator with a weighted probability
func (g *generator) Add(query RandomQuery) {
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
