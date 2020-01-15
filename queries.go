package main

import (
	"fmt"
	"io"
)

func genEthCall(w io.Writer, s State) error {
	r := s.RandInt64()
	fromBlock := s.CurrentBlock() - uint64(r%5000) // Pick a block within the last ~day
	toBlock := s.CurrentBlock() - uint64(r%5)      // Within the last ~minute
	// TODO: Use "latest" occasionally?
	address, topics := s.RandomContract()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_call","params":{"fromBlock":"0x%x","toBlock":"0x%x","address":"%s","topics":%s}}`+"\n", s.ID(), fromBlock, toBlock, address, topics)
	return err
}

func genEthGetTransactionReceipt(w io.Writer, s State) error {
	txID := s.RandomTransaction()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getTransactionReceipt","params":["%s"]}`+"\n", s.ID(), txID)
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
}
