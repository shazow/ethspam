package main

import (
	"context"
	"math/rand"
	"sync/atomic"

	"github.com/INFURA/go-ethlibs/eth"
	"github.com/INFURA/go-ethlibs/node"
)

type State interface {
	RandInt64() int64
	ID() int64
	CurrentBlock() int64
	RandomContract() (addr string, topics []string)
	RandomAddress() string
	RandomTransaction() string
}

type idGenerator struct {
	id int64
}

func (gen *idGenerator) Next() int64 {
	return atomic.AddInt64(&gen.id, 1)
}

// liveState implements State but it seeds the state dataset from live sources
// (Etherscan, etc)
type liveState struct {
	idGen   *idGenerator
	randSrc rand.Source

	currentBlock uint64
	transactions []string
}

func (s *liveState) ID() int64 {
	return s.idGen.Next()
}

func (s *liveState) CurrentBlock() uint64 {
	return s.currentBlock
}

func (s *liveState) RandInt64() int64 {
	return s.randSrc.Int63()
}

func (s *liveState) RandomTransaction() string {
	if len(s.transactions) == 0 {
		return ""
	}
	idx := int(s.randSrc.Int63()) % len(s.transactions)
	return s.transactions[idx]
}

func (s *liveState) RandomContract() (addr string, topics []string) {
	// TODO: Scrape https://etherscan.io/accounts instead?
	idx := s.RandInt64() % int64(len(popularContracts))
	c := popularContracts[idx]
	return c.Addr, c.Topics
}

var popularContracts = []struct {
	Addr   string
	Topics []string
}{
	{
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", // Wrapped ETH (wETH)
		[]string{"0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65", "0x00000000000000000000000076481caa104b5f6bccb540dae4cefaf1c398ebea"},
	},
	{
		"0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", // Unknown?
		[]string{"0x92ca3a80853e6663fa31fa10b99225f18d4902939b4c53a9caae9043f6efd004", "0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c"},
	},
	{
		"0x4ddc2d193948926d02f9b1fe9e1daa0718270ed5", // Compound ETH (cETH)
		[]string{"0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x0000000000000000000000004ddc2d193948926d02f9b1fe9e1daa0718270ed5"},
	},
}

type stateProducer struct {
	client node.Client
}

func (p *stateProducer) Refresh(oldState *liveState) (*liveState, error) {
	b, err := p.client.BlockByNumberOrTag(context.Background(), *(eth.MustBlockNumberOrTag("latest")), false)
	if err != nil {
		return nil, err
	}

	txs := make([]string, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		txs = append(txs, tx.Hash.String())
	}

	state := liveState{
		idGen:   oldState.idGen,
		randSrc: oldState.randSrc,

		currentBlock: b.Number.UInt64(),
		transactions: txs, // TODO: Keep some old transactions in the mix?
	}
	return &state, nil
}
