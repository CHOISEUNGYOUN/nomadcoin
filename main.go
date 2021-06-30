package main

import (
	"github.com/choiseungyoun/nomadcoin/explorer"
	"github.com/choiseungyoun/nomadcoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
