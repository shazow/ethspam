package main

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"

	"github.com/INFURA/go-ethlibs/eth"
	"github.com/INFURA/go-ethlibs/node"
)

var errEmptyBlock = errors.New("the sampled block is empty")

type State interface {
	RandInt64() int64
	ID() int64
	CurrentBlock() uint64
	RandomContract() (addr string, topics []string)
	RandomAddress() string
	RandomTransaction() string
	RandomCall() (to, from, input string, block uint64)
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
	transactions []eth.Transaction
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
	return s.transactions[idx].Hash.String()
}

func (s *liveState) RandomAddress() string {
	if len(s.transactions) == 0 {
		return ""
	}
	idx := int(s.randSrc.Int63()) % len(s.transactions)
	return s.transactions[idx].From.String()
}

func (s *liveState) RandomCall() (to, from, input string, block uint64) {
	if len(s.transactions) == 0 {
		return
	}
	tx := s.transactions[int(s.randSrc.Int63())%len(s.transactions)]
	if tx.To != nil {
		to = tx.To.String()
	}
	from = tx.From.String()
	input = tx.Input.String()
	block = tx.BlockNumber.UInt64()
	return
}

func (s *liveState) RandomContract() (addr string, topics []string) {
	// TODO: Scrape https://etherscan.io/accounts or https://ethgasstation.info/gasguzzlers.php instead?
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
	{
		"0x2a0c0dbecc7e4d658f48e01e3fa353f44050c208", // IDEX
		[]string{"0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7", "0xf341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567"},
	},
	{
		"0x06012c8cf97bead5deae237070f9587f8e7a266d", // CryptoKitties
		[]string{"0x241ea03ca20251805084d27d4440371c34a0b85ff108f6bb5611248f73818b80", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x0a5311bd2a6608f08a180df2ee7c5946819a649b204b554bb8e39825b2c50ad5"},
	},
	{
		"0x8d12a197cb00d4747a1fe03395095ce2a5cc6819", // EtherDelta
		[]string{"0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7", "0x1e0b760c386003e9cb9bcf4fcf3997886042859d9b6ed6320e804597fcdb28b0", "0xf341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567"},
	},
}

type stateProducer struct {
	client node.Client
}

func (p *stateProducer) Refresh(oldState *liveState) (*liveState, error) {
	if oldState == nil {
		return nil, errors.New("must provide old state to refresh")
	}

	b, err := p.client.BlockByNumberOrTag(context.Background(), *(eth.MustBlockNumberOrTag("latest")), true)
	if err != nil {
		return nil, err
	}
	// Short circuit if the sampled block is empty
	if len(b.Transactions) == 0 {
		return nil, errEmptyBlock
	}
	// txs will grow to the maximum contract transaction list size we'll see in a block, and the higher-indexed ones will stick around longer
	txs := oldState.transactions
	for i, tx := range b.Transactions {
		if tx.Transaction.Value.Int64() > 0 {
			// Only take 0-value transactions, hopefully these are all contract calls.
			continue
		}
		if len(oldState.transactions) < 50 || i > len(txs) {
			txs = append(txs, tx.Transaction)
			continue
		}
		// Keep some old transactions randomly
		if oldState.RandInt64()%6 > 2 {
			txs[i] = tx.Transaction
		}
	}

	state := liveState{
		idGen:   oldState.idGen,
		randSrc: oldState.randSrc,

		currentBlock: b.Number.UInt64(),
		transactions: txs,
	}
	return &state, nil
}
