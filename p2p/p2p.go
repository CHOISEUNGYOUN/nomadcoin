package p2p

import (
	"fmt"
	"net/http"

	"github.com/choiseungyoun/nomadcoin/blockchain"
	"github.com/choiseungyoun/nomadcoin/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	open_port := r.URL.Query().Get("open_port")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return open_port != "" || ip != ""
	}
	fmt.Printf("%s wants an upgrade\n", open_port)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, open_port)

}

func AddPeer(address, port, open_port string, broadcast bool) {
	fmt.Printf("%s to connect to port %s\n", open_port, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?open_port=%s", address, port, open_port), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	if broadcast {
		broadcaseNewPeer(p)
		return
	}
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func broadcaseNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
