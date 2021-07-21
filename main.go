package main

import (
	"github.com/choiseungyoun/nomadcoin/cli"
	"github.com/choiseungyoun/nomadcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
