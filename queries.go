package main

import (
	"fmt"
	"io"
	"strings"
)

// TODO: Replace with proper JSON serialization? Originally was written to be quick&dirty for maximum perf.

func genEthCall(w io.Writer, s State) error {
	// We eth_call the block before the call actually happened to avoid collision reverts
	to, from, input, block := s.RandomCall()
	var err error
	if to != "" {
		_, err = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_call","params":[{"to":%q,"from":%q,"data":%q},"0x%x"]}`+"\n", s.ID(), to, from, input, block-1)
	} else {
		_, err = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_call","params":[{"from":%q,"data":%q},"0x%x"]}`+"\n", s.ID(), from, input, block-1)
	}
	return err
}

func genEthGetTransactionReceipt(w io.Writer, s State) error {
	txID := s.RandomTransaction()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getTransactionReceipt","params":["%s"]}`+"\n", s.ID(), txID)
	return err
}

func genEthGetBalance(w io.Writer, s State) error {
	addr := s.RandomAddress()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getBalance","params":["%s","latest"]}`+"\n", s.ID(), addr)
	return err
}

func genEthGetBlockByNumber(w io.Writer, s State) error {
	r := s.RandInt64()
	// TODO: ~half of the block numbers are further from head
	blockNum := s.CurrentBlock() - uint64(r%5) // Within the last ~minute
	full := "true"
	if r%2 >= 0 {
		full = "false"
	}

	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getBlockByNumber","params":["0x%x",%s]}`+"\n", s.ID(), blockNum, full)
	return err
}

func genEthGetTransactionCount(w io.Writer, s State) error {
	addr := s.RandomAddress()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getTransactionCount","params":["%s","pending"]}`+"\n", s.ID(), addr)
	return err
}

func genEthBlockNumber(w io.Writer, s State) error {
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_blockNumber"}`+"\n", s.ID())
	return err
}

func genEthGetTransactionByHash(w io.Writer, s State) error {
	txID := s.RandomTransaction()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getTransactionByHash","params":["%s"]}`+"\n", s.ID(), txID)
	return err
}

func genEthGetLogs(w io.Writer, s State) error {
	r := s.RandInt64()
	// TODO: Favour latest/recent block on a curve
	fromBlock := s.CurrentBlock() - uint64(r%5000) // Pick a block within the last ~day
	toBlock := s.CurrentBlock() - uint64(r%5)      // Within the last ~minute
	address, topics := s.RandomContract()
	topicsJoined := strings.Join(topics, `","`)
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getLogs","params":[{"fromBlock":"0x%x","toBlock":"0x%x","address":"%s","topics":["%s"]}]}`+"\n", s.ID(), fromBlock, toBlock, address, topicsJoined)
	return err
}

func genEthGetCode(w io.Writer, s State) error {
	addr, _ := s.RandomContract()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getCode","params":["%s","latest"]}`+"\n", s.ID(), addr)
	return err
}

func installDefaults(gen *generator) {
	// Top queries by weight, pulled from a 5000 Infura query sample on Dec 2019.
	//     3 "eth_accounts"
	//     4 "eth_getStorageAt"
	//     4 "eth_syncing"
	//     7 "net_peerCount"
	//    12 "net_listening"
	//    14 "eth_gasPrice"
	//    16 "eth_sendRawTransaction"
	//    25 "net_version"
	//    30 "eth_getTransactionByBlockNumberAndIndex"
	//    38 "eth_getBlockByHash"
	//    45 "eth_estimateGas"
	//    88 "eth_getCode"
	//   252 "eth_getLogs"
	//   255 "eth_getTransactionByHash"
	//   333 "eth_blockNumber"
	//   390 "eth_getTransactionCount"
	//   399 "eth_getBlockByNumber"
	//   545 "eth_getBalance"
	//   607 "eth_getTransactionReceipt"
	//  1928 "eth_call"

	gen.Add(RandomQuery{
		Method:   "eth_call",
		Weight:   2000,
		Generate: genEthCall,
	})

	gen.Add(RandomQuery{
		Method:   "eth_getTransactionReceipt",
		Weight:   600,
		Generate: genEthGetTransactionReceipt,
	})

	gen.Add(RandomQuery{
		Method:   "eth_getBalance",
		Weight:   550,
		Generate: genEthGetBalance,
	})

	gen.Add(RandomQuery{
		Method:   "eth_getBlockByNumber",
		Weight:   400,
		Generate: genEthGetBalance,
	})

	gen.Add(RandomQuery{
		Method:   "eth_getTransactionCount",
		Weight:   400,
		Generate: genEthGetTransactionCount,
	})

	gen.Add(RandomQuery{
		Method:   "eth_blockNumber",
		Weight:   350,
		Generate: genEthBlockNumber,
	})

	gen.Add(RandomQuery{
		Method:   "eth_getTransactionByHash",
		Weight:   250,
		Generate: genEthGetTransactionByHash,
	})

	gen.Add(RandomQuery{
		Method:   "eth_getLogs",
		Weight:   250,
		Generate: genEthGetLogs,
	})

	gen.Add(RandomQuery{
		Method:   "eth_getCode",
		Weight:   100,
		Generate: genEthGetCode,
	})
}
