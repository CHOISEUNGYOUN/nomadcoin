package main

import (
	"github.com/choiseungyoun/nomadcoin/blockchain"
	"github.com/choiseungyoun/nomadcoin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
}
