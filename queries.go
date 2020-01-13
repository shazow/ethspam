package main

import (
	"fmt"
	"io"
)

func genEthCall(w io.Writer, s State) error {
	r := s.RandInt64()
	fromBlock := s.CurrentBlock() - (r % 5000) // Pick a block within the last ~day
	toBlock := s.CurrentBlock() - (r % 5)      // Within the last ~minute
	// TODO: Use "latest" occasionally?
	address, topics := s.RandomContract()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_call","params":{"fromBlock":"0x%x","toBlock":"0x%x","address":"%s","topics":%s}}`, s.ID(), fromBlock, toBlock, address, topics)
	return err
}

func genEthGetTransactionReceipt(w io.Writer, s State) error {
	txID := s.RandomTransaction()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getTransactionReceipt","params":["%s"]}`, s.ID(), txID)
	return err
}
