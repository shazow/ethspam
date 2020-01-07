package main

type State interface {
	RandInt64() int64
	ID() int
	Block(offset int64) int64
	RandomContract() (addr string, topics []string)
	RandomAddress() string
	RandomTransaction() string
}

// liveState implements State but it seeds the state dataset from live sources
// (Etherscan, etc)
type liveState struct {
}

func (s *liveState) RandomContract() (addr string, topics []string) {
	// TODO: Get some random popular contract
	return "0x931abd3732f7eada74190c8f89b46f8ba7103d54", []string{
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"0x0000000000000000000000000000000000000000000000000000000000000000",
	}
}
