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

func AddPeer(address, port, open_port string) {
	fmt.Printf("%s to connect to port %s\n", open_port, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?open_port=%s", address, port, open_port[1:]), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}
