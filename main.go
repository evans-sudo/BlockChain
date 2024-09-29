package main

import (
	"blockchainsystem/core"
	"blockchainsystem/crypto"
	"blockchainsystem/network"
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
)

//"blockchainsystem/network"

// server
// Transport ==> tcp, udp
// block
// Transacation
// key pairs,


func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemoteA := network.NewLocalTransport("REMOTE_A")
	trRemoteB := network.NewLocalTransport("REMOTE_B")
	trRemoteC := network.NewLocalTransport("REMOTE_C")


	trLocal.Connect(trRemoteA)
	trRemoteA.Connect(trRemoteB)
	trRemoteB.Connect(trRemoteC)
	trRemoteA.Connect(trLocal)


	initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

	go func ()  {
	for{
	//	trRemote.SendMessage(network.NetAddr(trLocal.Addr()), []byte("Hello world"))
		if err := SendTransaction(trRemoteA, network.NetAddr(trLocal.Addr())); err != nil {
			logrus.Error(err)
		}
		time.Sleep(2 * time.Second)
	}	
	}()


	privKey := crypto.GeneratePrivateKey()
	// opts := network.Serveropts{
	// 	PrivateKey: &privKey,
	// 	ID: "LOCAL",
	// 	Transport: []network.Transport{trLocal},
	// }

	localServer := makeServer("LOCAL", trLocal, &privKey)
	localServer.Start()
}

func initRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.Serveropts{
		PrivateKey: pk,
		ID: id,
		Transport: []network.Transport{tr},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}


func SendTransaction(tr network.Transport, to network.NetAddr) error {
	privkey := crypto.GeneratePrivateKey()
	data :=  []byte{0x02, 0x0a ,0x02, 0x0a, 0x0b}
	tx := core.NewTransaction(data)
	tx.Sign(privkey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err 
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}